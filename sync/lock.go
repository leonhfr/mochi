package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
)

const lockName = "mochi.lock"

type lockData map[string]lockDeck // indexed by deck id

type lockDeck struct {
	Path  string              `json:"path"`
	Name  string              `json:"name"`
	Cards map[string]lockCard `json:"cards,omitempty"` // indexed by card id
}

type lockCard struct {
	Filename string            `json:"filename"`         // filename inside directory: note.md
	Images   map[string]string `json:"images,omitempty"` // relative path: md5 hash
}

type Lock struct {
	data    lockData
	updated bool
	mu      sync.RWMutex
}

func ReadLock(ctx context.Context, client Client, fs filesystem.Interface) (*Lock, error) {
	source, err := fs.Read(lockName)
	if err != nil {
		return nil, err
	}
	if len(source) == 0 {
		return &Lock{data: lockData{}}, nil
	}

	data := lockData{}
	if err := json.Unmarshal(source, &data); err != nil {
		return nil, err
	}

	decks, err := client.ListDecks(ctx)
	if err != nil {
		return nil, err
	}

	lock := &Lock{data: data}
	lock.cleanDecks(decks)
	return lock, nil
}

func (l *Lock) Write(fs filesystem.Interface) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.updated {
		return false, nil
	}

	encoded, err := json.Marshal(l.data)
	if err != nil {
		return false, err
	}

	if err := fs.Write(lockName, string(encoded)); err != nil {
		return false, err
	}

	return true, nil
}

func (l *Lock) getDeck(path string) (string, lockDeck, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for deckID, deck := range l.data {
		if deck.Path == path {
			return deckID, deck, true
		}
	}

	return "", lockDeck{}, false
}

func (l *Lock) setDeck(id, path, name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.data[id] = lockDeck{Path: path, Name: name, Cards: map[string]lockCard{}}
	l.updated = true
}

func (l *Lock) getCard(deckID, cardID string) (lockCard, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, ok := l.data[deckID]; !ok {
		return lockCard{}, false
	}

	card, ok := l.data[deckID].Cards[cardID]
	return card, ok
}

func (l *Lock) setCard(deckID, cardID, filename string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.data[deckID]; !ok {
		return fmt.Errorf("deck %s not found", deckID)
	}

	if l.data[deckID].Cards == nil {
		l.data[deckID] = lockDeck{
			Path:  l.data[deckID].Path,
			Name:  l.data[deckID].Name,
			Cards: map[string]lockCard{},
		}
	}

	if _, ok := l.data[deckID].Cards[cardID]; ok {
		return nil
	}

	l.data[deckID].Cards[cardID] = lockCard{
		Filename: filename,
		Images:   map[string]string{},
	}
	l.updated = true
	return nil
}

func (l *Lock) getImageHash(deckID, cardID, path string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, ok := l.data[deckID]; !ok {
		return "", false
	}

	if _, ok := l.data[deckID].Cards[cardID]; !ok {
		return "", false
	}

	if l.data[deckID].Cards[cardID].Images == nil {
		l.data[deckID].Cards[cardID] = lockCard{
			Filename: l.data[deckID].Cards[cardID].Filename,
			Images:   make(map[string]string),
		}
	}

	hash, ok := l.data[deckID].Cards[cardID].Images[path]
	return hash, ok
}

func (l *Lock) setImageHash(deckID, cardID, path, hash string) error {
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

func (l *Lock) cleanDecks(apiDecks []api.Deck) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for deckID, deck := range l.data {
		index := slices.IndexFunc[[]api.Deck](apiDecks, func(apiDeck api.Deck) bool {
			return apiDeck.ID == deckID
		})

		if index < 0 {
			delete(l.data, deckID)
			l.updated = true
			continue
		}

		if apiDeck := apiDecks[index]; apiDeck.Name != deck.Name {
			deck.Name = apiDeck.Name
			l.data[deckID] = deck
			l.updated = true
		}
	}
}

func (l *Lock) cleanCards(deckID string, cardIDs []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.data[deckID]; !ok {
		return
	}

	for cardID := range l.data[deckID].Cards {
		if !slices.Contains[[]string](cardIDs, cardID) {
			delete(l.data[deckID].Cards, cardID)
			l.updated = true
		}
	}
}

func (l *Lock) cleanImages(deckID, cardID string, paths []string) {
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
