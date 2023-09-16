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
	Decks   map[string][2]string `toml:"decks,omitempty"` // directory path: [deck id, deck name]
	updated bool
	mu      sync.RWMutex
}

func ReadLock(ctx context.Context, client Client, fs filesystem.Interface) (*Lock, error) {
	source, err := fs.Read(lockName)
	if err != nil {
		return nil, err
	}

	lock := &Lock{}
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
	lock.mu.Lock()
	defer lock.mu.Unlock()

	for path, deck := range lock.Decks {
		index := slices.IndexFunc[[]api.Deck](decks, func(d api.Deck) bool {
			return deck[indexDeckID] == d.ID
		})

		if index < 0 {
			delete(lock.Decks, path)
			lock.updated = true
			continue
		}

		if apiDeck := decks[index]; deck[indexDeckName] != apiDeck.Name {
			lock.Decks[path] = [2]string{apiDeck.ID, apiDeck.Name}
			lock.updated = true
		}
	}
}
