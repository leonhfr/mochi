package lock

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"sync"

	"github.com/leonhfr/mochi/mochi"
)

const lockName = "mochi-lock.json"

type lockData map[string]Deck // indexed by deck id

// Deck contains the information about existing decks.
type Deck struct {
	ParentID string `json:"parentID"`
	Path     string `json:"path"`
	Name     string `json:"name"`
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
	lock := &Lock{path: path, rw: rw}
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

// GetDeck returns an existing decks information from a base string.
func (l *Lock) GetDeck(base string) (string, Deck, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for deckID, deck := range l.data {
		if deck.Path == base {
			return deckID, deck, true
		}
	}

	return "", Deck{}, false
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
