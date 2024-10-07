package deck

import (
	"context"
	"path/filepath"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// Client is the interface to interact with mochi decks.
type Client interface {
	CreateDeck(context.Context, mochi.CreateDeckRequest) (mochi.Deck, error)
	UpdateDeck(context.Context, string, mochi.UpdateDeckRequest) (mochi.Deck, error)
}

// Config is the interface to interact with the config.
type Config interface {
	GetDeck(string) (config.Deck, bool)
}

// Lockfile is the interface to interact with the lockfile.
type Lockfile interface {
	GetDeck(string) (string, lock.Deck, bool)
	SetDeck(string, string, string, string)
	UpdateDeckName(string, string)
}

// Sync syncs the name with the given path to mochi.
//
// It will create any intermediate decks as required until a root deck is reached.
// If the names do not match, the remote deck will be updated.
func Sync(ctx context.Context, client Client, config Config, lf Lockfile, path string) (deckID string, err error) {
	id, deck, ok := lf.GetDeck(path)
	if name := getDeckName(config, path); ok && deck.Name != name {
		err = updateDeckName(ctx, client, lf, deckID, name)
		return id, err
	} else if ok {
		return id, nil
	}

	parentID, stack := getStack(lf, path)
	for currentPath := ""; len(stack) > 0; {
		currentPath, stack = stack[len(stack)-1], stack[:len(stack)-1]
		name := getDeckName(config, currentPath)
		deckID, err = createDeck(ctx, client, lf, parentID, currentPath, name)
		if err != nil {
			return "", err
		}
		parentID = deckID
	}

	return
}

func createDeck(ctx context.Context, client Client, lf Lockfile, parentID, path, name string) (string, error) {
	deck, err := client.CreateDeck(ctx, mochi.CreateDeckRequest{
		Name:     name,
		ParentID: parentID,
	})
	if err != nil {
		return "", err
	}
	lf.SetDeck(deck.ID, parentID, path, name)
	return deck.ID, nil
}

func updateDeckName(ctx context.Context, client Client, lf Lockfile, deckID, name string) error {
	_, err := client.UpdateDeck(ctx, deckID, mochi.UpdateDeckRequest{Name: name})
	if err != nil {
		return err
	}
	lf.UpdateDeckName(deckID, name)
	return nil
}

func getDeckName(config Config, path string) string {
	deck, ok := config.GetDeck(path)
	if ok && deck.Name != nil {
		return *deck.Name
	}
	return filepath.Base(path)
}

func getStack(lockfile Lockfile, path string) (string, []string) {
	if path == "/" {
		return "", []string{path}
	}

	stack := []string{path}
	for path != "/" {
		path = filepath.Dir(path)
		deckID, _, ok := lockfile.GetDeck(path)
		if ok {
			return deckID, stack
		}
		stack = append(stack, path)
	}

	return "", stack
}
