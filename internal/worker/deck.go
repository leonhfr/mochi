package worker

import (
	"context"
	"slices"

	"github.com/sourcegraph/conc/iter"

	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/heap"
	"github.com/leonhfr/mochi/mochi"
)

// Deck contains a synced deck.
type Deck struct {
	deckID string
	cards  []deck.Card
}

// SyncDecks syncs the decks and parses the files.
func SyncDecks(ctx context.Context, logger Logger, r deck.Reader, p deck.Parser, client deck.CreateClient, config deck.CreateConfig, lf deck.CreateLockfile, workspace string, in <-chan heap.Group[heap.Path]) <-chan Result[Deck] {
	out := make(chan Result[Deck])
	go func() {
		defer close(out)
		for group := range in {
			deckConfig, ok := config.Deck(group.Base)
			if !ok {
				logger.Infof("parse(%s): discarding deck", group.Base)
				continue
			}

			logger.Infof("parse(%s): creating missing decks", group.Base)
			deckID, err := deck.Create(ctx, client, config, lf, group.Base)
			if err != nil {
				out <- Result[Deck]{err: err}
				continue
			}

			logger.Infof("parse(%s): parsing %d files", group.Base, len(group.Items))
			filePaths := heap.ConvertPaths(group.Items)
			cards, err := deck.Parse(r, p, workspace, deckConfig.Parser, filePaths)
			if err != nil {
				out <- Result[Deck]{err: err}
				continue
			}

			deckHeap := deck.Heap(cards)
			logger.Infof("parse(%s): parsed %d cards into %s decks", group.Base, len(cards), deckHeap.Len())
			for deckHeap.Len() > 0 {
				group := deckHeap.Pop()
				out <- Result[Deck]{
					data: Deck{
						deckID: deckID,
						cards:  group.Items,
					},
				}
			}
		}
	}()
	return out
}

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
