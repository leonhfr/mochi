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
	CreateDeck(ctx context.Context, req mochi.CreateDeckRequest) (mochi.Deck, error)
	UpdateDeck(ctx context.Context, id string, req mochi.UpdateDeckRequest) (mochi.Deck, error)
}

// Config is the interface to interact with the config.
type Config interface {
	GetDeck(base string) (config.Deck, bool)
}

// Lockfile is the interface to interact with the lockfile.
type Lockfile interface {
	GetDeck(base string) (string, lock.Deck, bool)
	SetDeck(id, parentID, path, name string)
	UpdateDeckName(id, name string)
}

// Sync sync the name with the given base to mochi.
//
// It will create any intermediate decks as required until a root deck is reached.
// If the names do not match, the remote deck will be updated.
func Sync(ctx context.Context, client Client, config Config, lf Lockfile, base string, name string) (deckID string, err error) {
	id, deck, ok := lf.GetDeck(base)
	if ok && name != "" && deck.Name != name {
		err = updateDeckName(ctx, client, lf, deckID, name)
		return id, err
	} else if ok {
		return id, nil
	}

	parentID, stack := getStack(lf, base)
	for currentBase := ""; len(stack) > 0; {
		currentBase, stack = stack[len(stack)-1], stack[:len(stack)-1]
		name := getDeckName(config, currentBase)
		deckID, err = createDeck(ctx, client, lf, parentID, currentBase, name)
		if err != nil {
			return "", err
		}
		parentID = deckID
	}

	return
}

func createDeck(ctx context.Context, client Client, lf Lockfile, parentID, base, name string) (string, error) {
	deck, err := client.CreateDeck(ctx, mochi.CreateDeckRequest{
		Name:     name,
		ParentID: parentID,
	})
	if err != nil {
		return "", err
	}
	lf.SetDeck(deck.ID, parentID, base, name)
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

func getDeckName(config Config, base string) string {
	deck, ok := config.GetDeck(base)
	if ok && len(deck.Name) > 0 {
		return deck.Name
	}
	return filepath.Base(base)
}

func getStack(lockfile Lockfile, base string) (string, []string) {
	if base == "/" {
		return "", []string{base}
	}

	stack := []string{base}
	for base != "/" {
		base = filepath.Dir(base)
		deckID, _, ok := lockfile.GetDeck(base)
		if ok {
			return deckID, stack
		}
		stack = append(stack, base)
	}

	return "", stack
}
