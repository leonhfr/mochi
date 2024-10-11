package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/mochi"
)

// Deck contains a synced deck.
type Deck struct {
	deckID    string
	config    config.Deck
	filePaths []string
}

// SyncDecks creates any missing decks and updates any mismatched name.
func SyncDecks(ctx context.Context, logger Logger, client deck.Client, config deck.Config, lf deck.Lockfile, in <-chan deck.Directory) <-chan Result[Deck] {
	out := make(chan Result[Deck], cap(in))
	go func() {
		defer close(out)

		for dir := range in {
			logger.Infof("deck sync: %s (%d files)", dir.Path, len(dir.FilePaths))
			deckConfig, ok := config.GetDeck(dir.Path)
			if !ok {
				logger.Infof("deck sync: discarded %s", dir.Path)
				continue
			}

			deckID, err := deck.Sync(ctx, client, config, lf, dir.Path)
			logger.Infof("deck sync(deckID %s): synced %s", deckID, dir.Path)

			out <- Result[Deck]{
				data: Deck{
					deckID:    deckID,
					config:    deckConfig,
					filePaths: dir.FilePaths,
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
