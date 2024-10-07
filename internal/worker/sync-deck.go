package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// SyncedDeck contains a synced deck.
type SyncedDeck struct {
	id        string
	filePaths []string
}

// SyncDeck creates any missing decks and updates any mismatched name.
func SyncDeck(ctx context.Context, logger Logger, client *mochi.Client, cfg *config.Config, lf *lock.Lock, in <-chan FilteredDeck) <-chan Result[SyncedDeck] {
	out := make(chan Result[SyncedDeck])
	go func() {
		defer close(out)

		for filteredDeck := range in {
			logger.Infof("deck sync: %s", filteredDeck.path)
			deckID, err := deck.Sync(ctx, client, cfg, lf, filteredDeck.path)
			out <- Result[SyncedDeck]{
				data: SyncedDeck{id: deckID, filePaths: filteredDeck.filePaths},
				err:  err,
			}
		}
	}()
	return out
}
