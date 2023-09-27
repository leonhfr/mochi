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

func SynchronizeDecks(ctx context.Context, sources []string, lock *Lock, config Config, client Client, logger Logger) (DeckResult, error) {
	var res DeckResult
	for _, path := range uniqueDirs(sources) {
		name := config.deckName(path)
		deckID, deck, ok := lock.getDeck(path)

		if !ok {
			if err := createDeck(ctx, path, name, lock, client); err != nil {
				return res, err
			}
			logger.Infof("Created deck \"%s\"", name)
			res.Created++
		}

		if ok && deck.Name != name {
			if err := updateDeck(ctx, path, deckID, name, lock, client); err != nil {
				return res, err
			}
			logger.Infof("Updated deck \"%s\"", name)
			res.Updated++
		}
	}
	return res, nil
}

func createDeck(ctx context.Context, path, name string, lock *Lock, client Client) error {
	var parentID string
	if path := filepath.Dir(path); len(path) > 1 {
		if deckID, _, ok := lock.getDeck(path); ok {
			parentID = deckID
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

	lock.setDeck(deck.ID, path, deck.Name)
	return nil
}

func updateDeck(ctx context.Context, path, id, name string, lock *Lock, client Client) error {
	req := api.UpdateDeckRequest{
		Name: name,
	}

	if _, err := client.UpdateDeck(ctx, id, req); err != nil {
		return err
	}

	deckID, _, _ := lock.getDeck(path)
	lock.setDeck(deckID, path, name)
	return nil
}
