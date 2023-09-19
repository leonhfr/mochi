package sync

import (
	"bytes"
	"context"
	"sync"

	"github.com/BurntSushi/toml"
	"golang.org/x/exp/slices"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
)

const (
	lockName   = "mochi-lock.toml"
	lockIndent = ""
)

const (
	indexDeckID int = iota
	indexDeckName
)

type Lock struct {
	Decks   map[string][2]string                    `toml:"decks,omitempty"`  // directory path: [deck id, deck name]
	Images  map[string]map[string]map[string]string `toml:"images,omitempty"` // deck id: card id: file path: hash
	updated bool
	mu      sync.RWMutex
}

func ReadLock(ctx context.Context, client Client, fs filesystem.Interface) (*Lock, error) {
	source, err := fs.Read(lockName)
	if err != nil {
		return nil, err
	}

	lock := &Lock{
		Decks:  map[string][2]string{},
		Images: map[string]map[string]map[string]string{},
	}
	if err := toml.Unmarshal(source, lock); err != nil {
		return nil, err
	}

	decks, err := client.ListDecks(ctx)
	if err != nil {
		return nil, err
	}

	updateLock(lock, decks)

	return lock, nil
}

func (l *Lock) getDeck(path string) ([2]string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	deck, ok := l.Decks[path]
	return deck, ok
}

func (l *Lock) setDeck(path string, deck [2]string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Decks[path] = deck
	l.updated = true
}

func (l *Lock) deleteDeck(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.Decks, path)
	l.updated = true
}

func (l *Lock) getImageCards(deckID string) []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	deck, ok := l.Images[deckID]
	if !ok {
		return nil
	}

	cards := make([]string, 0, len(deck))
	for card := range deck {
		cards = append(cards, card)
	}
	return cards
}

func (l *Lock) deleteImageDeck(deckID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.Images, deckID)
	l.updated = true
}

func (l *Lock) deleteImageCard(deckID, cardID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.Images[deckID], cardID)
	l.updated = true
}

func (l *Lock) getImageHash(deckID, cardID, path string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	deckImageHashes, ok := l.Images[deckID]
	if !ok {
		return "", false
	}

	cardImageHashes, ok := deckImageHashes[cardID]
	if !ok {
		return "", false
	}

	hash, ok := cardImageHashes[path]
	return hash, ok
}

func (l *Lock) setImageHash(deckID, cardID, path, hash string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.updated = true
	if _, ok := l.Images[deckID]; !ok {
		l.Images[deckID] = map[string]map[string]string{
			cardID: {path: hash},
		}
		return
	}

	if _, ok := l.Images[deckID][cardID]; !ok {
		l.Images[deckID][cardID] = map[string]string{
			path: hash,
		}
		return
	}

	l.Images[deckID][cardID][path] = hash
}

func (l *Lock) Write(fs filesystem.Interface) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.updated {
		return false, nil
	}

	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Indent = lockIndent
	if err := encoder.Encode(l); err != nil {
		return false, err
	}

	if err := fs.Write(lockName, buf.String()); err != nil {
		return false, err
	}

	return true, nil
}

func updateLock(lock *Lock, decks []api.Deck) {
	for path, deck := range lock.Decks {
		index := slices.IndexFunc[[]api.Deck](decks, func(d api.Deck) bool {
			return deck[indexDeckID] == d.ID
		})

		if index < 0 {
			lock.deleteDeck(path)
			continue
		}

		if apiDeck := decks[index]; deck[indexDeckName] != apiDeck.Name {
			lock.setDeck(path, [2]string{apiDeck.ID, apiDeck.Name})
		}
	}

	for deckID := range lock.Images {
		if !slices.ContainsFunc[[]api.Deck](decks, func(d api.Deck) bool {
			return d.ID == deckID
		}) {
			lock.deleteImageDeck(deckID)
		}
	}
}
