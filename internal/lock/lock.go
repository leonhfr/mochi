package lock

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"sync"

	"github.com/go-playground/validator/v10"

	"github.com/leonhfr/mochi/mochi"
)

const lockName = "mochi-lock.json"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Lock represents a lockfile.
type Lock struct {
	decks   map[string]Deck // indexed by deck id
	path    string
	updated bool
	mu      sync.Mutex
	rw      ReaderWriter
}

// Deck contains the information about existing decks.
type Deck struct {
	ParentID string          `json:"parentID,omitempty"`
	Path     string          `json:"path" validate:"required"`
	Name     string          `json:"name" validate:"required"`
	Cards    map[string]Card `json:"cards,omitempty" validate:"dive"` // indexed by card id
}

// Card contains the information about existing cards.
type Card struct {
	Filename string            `json:"filename" validate:"required"` // filename inside directory: note.md
	Images   map[string]string `json:"images,omitempty"`             // map[path]md5 hash
}

// ReaderWriter represents the interface to interact with a lockfile.
type ReaderWriter interface {
	Read(string) (io.ReadCloser, error)
	Write(string) (io.WriteCloser, error)
}

// Parse parses the lockfile in the target directory.
func Parse(rw ReaderWriter, target string) (*Lock, error) {
	path := filepath.Join(target, lockName)
	lock := &Lock{
		decks: make(map[string]Deck),
		path:  path,
		rw:    rw,
	}

	r, err := rw.Read(path)
	if err == fs.ErrNotExist {
		return lock, nil
	} else if err != nil {
		return nil, err
	}
	defer r.Close()

	if err := json.NewDecoder(r).Decode(&lock.decks); err != nil {
		return nil, err
	}

	for _, data := range lock.decks {
		if err := validate.Struct(&data); err != nil {
			return nil, err
		}
	}

	return lock, nil
}

// Lock locks the lockfile.
func (l *Lock) Lock() {
	l.mu.Lock()
}

// Unlock unlocks the lockfile.
func (l *Lock) Unlock() {
	l.mu.Unlock()
}

// CleanDecks removes from the lockfile the inexistent decks
// and updates the deck names if they differ.
func (l *Lock) CleanDecks(decks []mochi.Deck) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for id, lockDeck := range l.decks {
		index := slices.IndexFunc(decks, func(deck mochi.Deck) bool {
			return deck.ID == id
		})

		if index < 0 {
			delete(l.decks, id)
			l.updated = true
			continue
		}

		if decks[index].ParentID != lockDeck.ParentID {
			delete(l.decks, id)
			l.updated = true
			continue
		}

		if decks[index].Name != lockDeck.Name {
			lockDeck.Name = decks[index].Name
			l.decks[id] = lockDeck
			l.updated = true
		}
	}
}

// CleanCards removes from the lockfile the inexistent cards in a deck.
func (l *Lock) CleanCards(deckID string, cardIDs []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.decks[deckID]; !ok {
		return
	}

	for cardID := range l.decks[deckID].Cards {
		if !slices.Contains(cardIDs, cardID) {
			delete(l.decks[deckID].Cards, cardID)
			l.updated = true
		}
	}
}

// CleanImages removes from the lockfile the inexistent paths in a card.
func (l *Lock) CleanImages(deckID, cardID string, paths []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.decks[deckID]; !ok {
		return
	}

	if _, ok := l.decks[deckID].Cards[cardID]; !ok {
		return
	}

	for path := range l.decks[deckID].Cards[cardID].Images {
		if !slices.Contains[[]string](paths, path) {
			delete(l.decks[deckID].Cards[cardID].Images, path)
			l.updated = true
		}
	}
}

// Deck returns an existing decks information from a directory string.
//
// Assumes mutex is already acquired.
func (l *Lock) Deck(path string) (string, Deck, bool) {
	for deckID, deck := range l.decks {
		if deck.Path == path {
			return deckID, deck, true
		}
	}

	return "", Deck{}, false
}

// SetDeck sets a deck in the lockfile.
//
// Assumes mutex is already acquired.
func (l *Lock) SetDeck(id, parentID, path, name string) {
	l.decks[id] = Deck{
		ParentID: parentID,
		Path:     path,
		Name:     name,
		Cards:    make(map[string]Card),
	}
	l.updated = true
}

// UpdateDeckName updates a deck name in the lockfile.
//
// Assumes mutex is already acquired.
func (l *Lock) UpdateDeckName(id, name string) {
	tmp := l.decks[id]
	tmp.Name = name
	l.decks[id] = tmp
	l.updated = true
}

// Card returns an existing cards data.
func (l *Lock) Card(deckID, cardID string) (Card, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.decks[deckID]; !ok {
		return Card{}, false
	}

	card, ok := l.decks[deckID].Cards[cardID]
	return card, ok
}

// SetCard sets a card in the given deck.
func (l *Lock) SetCard(deckID, cardID, filename string, images map[string]string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.decks[deckID]; !ok {
		return fmt.Errorf("deck %s not found", deckID)
	}

	if l.decks[deckID].Cards == nil {
		l.decks[deckID] = Deck{
			Path:  l.decks[deckID].Path,
			Name:  l.decks[deckID].Name,
			Cards: map[string]Card{},
		}
	}

	l.decks[deckID].Cards[cardID] = Card{
		Filename: filename,
		Images:   images,
	}
	l.updated = true

	return nil
}

// ImageHashes returns the image hashes.
// If the image does not exist an empty string is returned.
func (l *Lock) ImageHashes(deckID, cardID string, paths []string) []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	hashes := make([]string, 0, len(paths))
	for _, path := range paths {
		hashes = append(hashes, l.getImageHash(deckID, cardID, path))
	}

	return hashes
}

// requires read lock to be already acquired.
func (l *Lock) getImageHash(deckID, cardID, path string) string {
	if _, ok := l.decks[deckID]; !ok {
		return ""
	}

	if _, ok := l.decks[deckID].Cards[cardID]; !ok {
		return ""
	}

	return l.decks[deckID].Cards[cardID].Images[path]
}

// Updated returns whether the lockfile has been updated.
func (l *Lock) Updated() bool {
	return l.updated
}

// String implements fmt.Stringer.
func (l *Lock) String() string {
	return fmt.Sprintf("data(updated: %t): %v", l.updated, l.decks)
}

// Write writes the lockfile to the target directory.
func (l *Lock) Write() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.updated {
		return nil
	}

	w, err := l.rw.Write(l.path)
	if err != nil {
		return err
	}
	defer w.Close()

	return json.NewEncoder(w).Encode(l.decks)
}
