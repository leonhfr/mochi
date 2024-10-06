package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// DeckSync creates any missing decks and
// updates any mismatched name.
func DeckSync(ctx context.Context, logger Logger, client *mochi.Client, cfg *config.Config, lf *lock.Lock, in <-chan DeckJob) <-chan Result[DeckJob] {
	out := make(chan Result[DeckJob])
	go func() {
		defer close(out)

		for d := range in {
			logger.Infof("deck sync: %s", d.dir.Path)
			deckID, err := deck.Sync(ctx, client, cfg, lf, d.dir.Path, d.cfg.Name)
			d.id = deckID
			out <- Result[DeckJob]{
				data: d,
				err:  err,
			}
		}
	}()
	return out
}
