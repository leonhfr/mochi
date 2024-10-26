package worker

import (
	"context"
	"slices"

	"github.com/sourcegraph/conc/iter"

	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/mochi"
)

// DeckListClient is the interface to implement to list decks.
type DeckListClient interface {
	ListDecks(ctx context.Context) ([]mochi.Deck, error)
}

// ListDecks lists existing decks.
func ListDecks(ctx context.Context, client DeckListClient) (<-chan string, error) {
	decks, err := client.ListDecks(ctx)
	if err != nil {
		return nil, err
	}
	out := make(chan string)
	go func() {
		defer close(out)
		for _, deck := range decks {
			select {
			case out <- deck.ID:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

// CleanDecksClient is the client interface to clean decks.
type CleanDecksClient interface {
	DeckListClient
	deck.DeleteEmptyClient
}

// DeleteLeafDecks cleans the decks and returns true if at least one deck has been cleaned.
func DeleteLeafDecks(ctx context.Context, client CleanDecksClient) (bool, error) {
	decks, err := client.ListDecks(ctx)
	if err != nil {
		return false, err
	}

	leaves := deck.LeafDecks(decks)
	if len(leaves) == 0 {
		return false, nil
	}

	done, err := iter.MapErr(leaves, func(deckID *string) (bool, error) {
		return deck.DeleteEmpty(ctx, client, *deckID)
	})

	return slices.Contains(done, true), err
}
