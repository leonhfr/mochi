package lock

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/go-playground/validator/v10"
)

const lockName = "mochi-lock.json"

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		deck := sl.Current().Interface().(Deck)
		if len(deck.Path) == 0 && !deck.Virtual {
			sl.ReportError(deck.Path, "path", "Path", "pathorvirtual", "")
			sl.ReportError(deck.Virtual, "virtual", "Virtual", "pathorvirtual", "")
		}
		if len(deck.Path) > 0 && deck.Virtual {
			sl.ReportError(deck.Path, "path", "Path", "pathorvirtual", "")
			sl.ReportError(deck.Virtual, "virtual", "Virtual", "pathorvirtual", "")
		}
	}, Deck{})
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
	Path     string          `json:"path,omitempty"`
	Name     string          `json:"name" validate:"required"`
	Cards    map[string]Card `json:"cards,omitempty" validate:"dive"` // indexed by card id
	Virtual  bool            `json:"virtual,omitempty"`
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

// Decks returns all decks.
func (l *Lock) Decks() map[string]Deck {
	return l.decks
}

// VirtualDecks returns all virtual decks for a parent deck.
func (l *Lock) VirtualDecks(parentID string) map[string]Deck {
	decks := make(map[string]Deck)
	for id, deck := range l.decks {
		if deck.ParentID == parentID && deck.Virtual {
			decks[id] = deck
		}
	}
	return decks
}

// Deck returns a deck.
//
// Assumes mutex is already acquired.
func (l *Lock) Deck(id string) (Deck, bool) {
	deck, ok := l.decks[id]
	return deck, ok
}

// DeckFromPath returns an existing decks information from a directory string.
//
// Assumes mutex is already acquired.
func (l *Lock) DeckFromPath(path string) (string, Deck, bool) {
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

// SetVirtualDeck sets a virtual deck in the lockfile.
//
// Assumes mutex is already acquired.
func (l *Lock) SetVirtualDeck(id, parentID, name string) {
	l.decks[id] = Deck{
		ParentID: parentID,
		Name:     name,
		Cards:    make(map[string]Card),
		Virtual:  true,
	}
	l.updated = true
}

// UpdateDeck updates a deck name in the lockfile.
//
// Assumes mutex is already acquired.
func (l *Lock) UpdateDeck(id, name string) {
	tmp := l.decks[id]
	tmp.Name = name
	l.decks[id] = tmp
	l.updated = true
}

// DeleteDeck deletes a deck from the lockfile.
//
// Assumes mutex is already acquired.
func (l *Lock) DeleteDeck(id string) {
	delete(l.decks, id)
	l.updated = true
}

// Card returns an existing cards data.
//
// Assumes mutex is already acquired.
func (l *Lock) Card(deckID, cardID string) (Card, bool) {
	if _, ok := l.decks[deckID]; !ok {
		return Card{}, false
	}

	card, ok := l.decks[deckID].Cards[cardID]
	return card, ok
}

// SetCard sets a card in the given deck.
//
// Assumes mutex is already acquired.
func (l *Lock) SetCard(deckID, cardID, filename string, images map[string]string) error {
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

// DeleteCard deletes a card in the given deck.
//
// Assumes mutex is already acquired.
func (l *Lock) DeleteCard(deckID, cardID string) {
	delete(l.decks[deckID].Cards, cardID)
	l.updated = true
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
