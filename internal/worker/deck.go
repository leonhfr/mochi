package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// FilteredDeck represents a deck whose config has been matched.
type FilteredDeck struct {
	path      string
	filePaths []string
}

// FilterDeck filters the directories, only forwarding them
// if a deck config has been found.
func FilterDeck(logger Logger, cfg *config.Config, in <-chan deck.Directory) <-chan FilteredDeck {
	out := make(chan FilteredDeck)
	go func() {
		defer close(out)

		for dir := range in {
			if deckConfig, ok := cfg.GetDeck(dir.Path); ok {
				logger.Infof("deck filter: forwarding %s with %d files", dir.Path, len(dir.FilePaths))
				out <- FilteredDeck{
					path:      deckConfig.Path,
					filePaths: dir.FilePaths,
				}
			} else {
				logger.Infof("deck filter: discarded %s with %d files", dir.Path, len(dir.FilePaths))
			}
		}
	}()
	return out
}

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
