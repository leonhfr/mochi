package sync

import (
	"context"
	"path/filepath"

	"github.com/leonhfr/mochi/api"
)

const rootDeckName = "Root"

type DeckResult struct {
	Created int
	Updated int
}

func SynchronizeDecks(ctx context.Context, sources []string, lock *Lock, config Config, client Client) (DeckResult, error) {
	var res DeckResult
	for _, path := range uniqueDirs(sources) {
		name := config.deckName(path)
		deck, ok := lock.getDeck(path)

		if !ok {
			if err := createDeck(ctx, path, name, lock, client); err != nil {
				return res, err
			}
			res.Created++
		}

		if ok && deck[indexDeckName] != name {
			if err := updateDeck(ctx, path, deck[indexDeckID], name, lock, client); err != nil {
				return res, err
			}
			res.Updated++
		}
	}
	return res, nil
}

func createDeck(ctx context.Context, path, name string, lock *Lock, client Client) error {
	var parentID string
	if path := filepath.Dir(path); len(path) > 1 {
		if deck, ok := lock.getDeck(path); ok {
			parentID = deck[indexDeckID]
		}
	}

	req := api.CreateDeckRequest{
		Name:     name,
		ParentID: parentID,
	}

	deck, err := client.CreateDeck(ctx, req)
	if err != nil {
		return err
	}

	lock.setDeck(path, [2]string{deck.ID, deck.Name})
	return nil
}

func updateDeck(ctx context.Context, path, id, name string, lock *Lock, client Client) error {
	req := api.UpdateDeckRequest{
		Name: name,
	}

	if _, err := client.UpdateDeck(ctx, id, req); err != nil {
		return err
	}

	deck, _ := lock.getDeck(path)
	lock.setDeck(path, [2]string{deck[indexDeckID], name})
	return nil
}
