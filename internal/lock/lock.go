package lock

import (
	"encoding/json"
	"fmt"
	"io"
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

type lockData map[string]Deck // indexed by deck id

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

// Lock represents a lockfile.
type Lock struct {
	data    lockData
	path    string
	updated bool
	mu      sync.RWMutex
	rw      ReaderWriter
}

// ReaderWriter represents the interface to interact with a lockfile.
type ReaderWriter interface {
	Exists(string) bool
	Read(string) (io.ReadCloser, error)
	Write(string) (io.WriteCloser, error)
}

// Parse parses the lockfile in the target directory.
func Parse(rw ReaderWriter, target string) (*Lock, error) {
	path := filepath.Join(target, lockName)
	lock := &Lock{
		data: make(lockData),
		path: path,
		rw:   rw,
	}
	if !rw.Exists(path) {
		return lock, nil
	}

	r, err := rw.Read(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if err := json.NewDecoder(r).Decode(&lock.data); err != nil {
		return nil, err
	}

	for _, data := range lock.data {
		if err := validate.Struct(&data); err != nil {
			return nil, err
		}
	}

	return lock, nil
}

// CleanDecks removes from the lockfile the inexistent decks
// and updates the deck names if they differ.
func (l *Lock) CleanDecks(decks []mochi.Deck) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for id, lockDeck := range l.data {
		index := slices.IndexFunc(decks, func(deck mochi.Deck) bool {
			return deck.ID == id
		})

		if index < 0 {
			delete(l.data, id)
			l.updated = true
			continue
		}

		if decks[index].ParentID != lockDeck.ParentID {
			delete(l.data, id)
			l.updated = true
			continue
		}

		if decks[index].Name != lockDeck.Name {
			lockDeck.Name = decks[index].Name
			l.data[id] = lockDeck
			l.updated = true
		}
	}
}

// CleanCards removes from the lockfile the inexistent cards in a deck.
func (l *Lock) CleanCards(deckID string, cardIDs []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.data[deckID]; !ok {
		return
	}

	for cardID := range l.data[deckID].Cards {
		if !slices.Contains(cardIDs, cardID) {
			delete(l.data[deckID].Cards, cardID)
			l.updated = true
		}
	}
}

// CleanImages removes from the lockfile the inexistent paths in a card.
func (l *Lock) CleanImages(deckID, cardID string, paths []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.data[deckID]; !ok {
		return
	}

	if _, ok := l.data[deckID].Cards[cardID]; !ok {
		return
	}

	for path := range l.data[deckID].Cards[cardID].Images {
		if !slices.Contains[[]string](paths, path) {
			delete(l.data[deckID].Cards[cardID].Images, path)
			l.updated = true
		}
	}
}

// GetDeck returns an existing decks information from a directory string.
func (l *Lock) GetDeck(path string) (string, Deck, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for deckID, deck := range l.data {
		if deck.Path == path {
			return deckID, deck, true
		}
	}

	return "", Deck{}, false
}

// SetDeck sets a deck in the lockfile.
func (l *Lock) SetDeck(id, parentID, path, name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.data[id] = Deck{
		ParentID: parentID,
		Path:     path,
		Name:     name,
		Cards:    make(map[string]Card),
	}
	l.updated = true
}

// UpdateDeckName updates a deck name in the lockfile.
func (l *Lock) UpdateDeckName(id, name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	tmp := l.data[id]
	tmp.Name = name
	l.data[id] = tmp
	l.updated = true
}

// GetCard returns an existing cards data.
func (l *Lock) GetCard(deckID string, cardID string) (Card, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, ok := l.data[deckID]; !ok {
		return Card{}, false
	}

	card, ok := l.data[deckID].Cards[cardID]
	return card, ok
}

// SetCard sets a card in the given deck.
func (l *Lock) SetCard(deckID string, cardID string, filename string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.data[deckID]; !ok {
		return fmt.Errorf("deck %s not found", deckID)
	}

	if l.data[deckID].Cards == nil {
		l.data[deckID] = Deck{
			Path:  l.data[deckID].Path,
			Name:  l.data[deckID].Name,
			Cards: map[string]Card{},
		}
	}

	if _, ok := l.data[deckID].Cards[cardID]; ok {
		return nil
	}

	l.data[deckID].Cards[cardID] = Card{
		Filename: filename,
	}
	l.updated = true
	return nil
}

// GetImageHash returns the image hash and ok if it exists.
func (l *Lock) GetImageHash(deckID, cardID, path string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, ok := l.data[deckID]; !ok {
		return "", false
	}

	if _, ok := l.data[deckID].Cards[cardID]; !ok {
		return "", false
	}

	if l.data[deckID].Cards[cardID].Images == nil {
		l.data[deckID].Cards[cardID] = Card{
			Filename: l.data[deckID].Cards[cardID].Filename,
			Images:   make(map[string]string),
		}
	}

	hash, ok := l.data[deckID].Cards[cardID].Images[path]
	return hash, ok
}

// SetImageHash sets the hash to an image path.
func (l *Lock) SetImageHash(deckID, cardID, path, hash string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.data[deckID]; !ok {
		return fmt.Errorf("deck %s not found", deckID)
	}

	if _, ok := l.data[deckID].Cards[cardID]; !ok {
		return fmt.Errorf("card %s not found in deck %s", cardID, deckID)
	}

	l.data[deckID].Cards[cardID].Images[path] = hash
	l.updated = true
	return nil
}

// Updated returns whether the lockfile has been updated.
func (l *Lock) Updated() bool {
	return l.updated
}

// String implements fmt.Stringer.
func (l *Lock) String() string {
	return fmt.Sprintf("data(updated: %t): %v", l.updated, l.data)
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

	return json.NewEncoder(w).Encode(l.data)
}
