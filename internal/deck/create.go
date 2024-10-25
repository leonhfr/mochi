package deck

import (
	"context"
	"path/filepath"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// CreateClient is the interface to sync mochi decks.
type CreateClient interface {
	CreateDeck(context.Context, mochi.CreateDeckRequest) (mochi.Deck, error)
	UpdateDeck(context.Context, string, mochi.UpdateDeckRequest) (mochi.Deck, error)
}

// CreateConfig is the interface to interact with the config.
type CreateConfig interface {
	Deck(string) (config.Deck, bool)
}

// CreateLockfile is the interface to interact with the lockfile.
type CreateLockfile interface {
	Lock()
	Unlock()
	DeckFromPath(string) (string, lock.Deck, bool)
	SetDeck(string, string, string, string)
	UpdateDeck(string, string)
}

// Create creates the deck.
//
// It will create any intermediate decks as required until a root deck is reached.
// If the names do not match, the remote deck will be updated.
func Create(ctx context.Context, client CreateClient, config CreateConfig, lf CreateLockfile, path string) (deckID string, err error) {
	lf.Lock()
	defer lf.Unlock()

	id, deck, ok := lf.DeckFromPath(path)
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

func createDeck(ctx context.Context, client CreateClient, lf CreateLockfile, parentID, path, name string) (string, error) {
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

func updateDeckName(ctx context.Context, client CreateClient, lf CreateLockfile, deckID, name string) error {
	_, err := client.UpdateDeck(ctx, deckID, mochi.UpdateDeckRequest{Name: name})
	if err != nil {
		return err
	}
	lf.UpdateDeck(deckID, name)
	return nil
}

func getDeckName(config CreateConfig, path string) string {
	deck, ok := config.Deck(path)
	if ok && deck.Path == path && deck.Name != "" {
		return deck.Name
	}
	return filepath.Base(path)
}

func getStack(lockfile CreateLockfile, path string) (string, []string) {
	if path == "/" {
		return "", []string{path}
	}

	stack := []string{path}
	for !isTopLevelDirectory(path) {
		path = filepath.Dir(path)
		deckID, _, ok := lockfile.DeckFromPath(path)
		if ok {
			return deckID, stack
		}
		stack = append(stack, path)
	}

	return "", stack
}

func isTopLevelDirectory(path string) bool {
	dir, _ := filepath.Split(path)
	return dir == "/"
}
