package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// Deck contains a synced deck.
type Deck struct {
	deckID    string
	config    config.Deck
	filePaths []string
}

// SyncDecks creates any missing decks and updates any mismatched name.
func SyncDecks(ctx context.Context, logger Logger, client *mochi.Client, config *config.Config, lf *lock.Lock, in <-chan deck.Directory) <-chan Result[Deck] {
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
