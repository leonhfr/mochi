package worker

import (
	"context"
	"slices"

	"github.com/sourcegraph/conc/iter"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/sync"
	"github.com/leonhfr/mochi/mochi"
)

// Deck contains a synced deck.
type Deck struct {
	deckID    string
	config    config.Deck
	filePaths []string
}

// SyncDecks creates any missing decks and updates any mismatched name.
func SyncDecks(ctx context.Context, logger Logger, client deck.Client, config deck.Config, lf deck.Lockfile, in <-chan sync.Group[sync.Path]) <-chan Result[Deck] {
	out := make(chan Result[Deck], cap(in))
	go func() {
		defer close(out)

		for dir := range in {
			logger.Infof("deck sync: %s (%d files)", dir.Base, len(dir.Items))
			deckConfig, ok := config.GetDeck(dir.Base)
			if !ok {
				logger.Infof("deck sync: discarded %s", dir.Base)
				continue
			}

			deckID, err := deck.Sync(ctx, client, config, lf, dir.Base)
			logger.Infof("deck sync(deckID %s): synced %s", deckID, dir.Base)

			out <- Result[Deck]{
				data: Deck{
					deckID:    deckID,
					config:    deckConfig,
					filePaths: sync.ConvertPaths(dir.Items),
				},
				err: err,
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
	out := make(chan string, len(decks))
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
	deck.CleanClient
}

// CleanDecks cleans the decks and returns true if at least one deck has been cleaned.
func CleanDecks(ctx context.Context, client CleanDecksClient) (bool, error) {
	decks, err := client.ListDecks(ctx)
	if err != nil {
		return false, err
	}

	leaves := deck.LeafDecks(decks)
	if len(leaves) == 0 {
		return false, nil
	}

	done, err := iter.MapErr(leaves, func(deckID *string) (bool, error) {
		return deck.Clean(ctx, client, *deckID)
	})

	return slices.Contains(done, true), err
}
