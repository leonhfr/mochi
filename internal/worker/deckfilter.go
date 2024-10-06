package worker

import (
	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
)

// DeckJob contains a deck config and its associated files.
type DeckJob struct {
	dir deck.Directory
	cfg config.Deck
}

// DeckFilter filters the directories, only forwarding them
// if a deck config has been found.
func DeckFilter(logger Logger, cfg *config.Config, in <-chan deck.Directory) <-chan DeckJob {
	out := make(chan DeckJob)
	go func() {
		defer close(out)

		for dir := range in {
			if deck, ok := cfg.GetDeck(dir.Base); ok {
				logger.Infof("deck filter: forwarding %s with %d files", dir.Base, len(dir.Paths))
				out <- DeckJob{dir: dir, cfg: deck}
			} else {
				logger.Infof("deck filter: discarded %s with %d files", dir.Base, len(dir.Paths))
			}
		}
	}()
	return out
}
