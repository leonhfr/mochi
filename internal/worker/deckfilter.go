package worker

import (
	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
)

// Deck contains a deck config and its associated files.
type Deck struct {
	dir deck.Directory
	cfg config.Deck
}

// DeckFilter filters the directories, only forwarding them
// if a deck config has been found.
func DeckFilter(logger Logger, cfg *config.Config, in <-chan deck.Directory) <-chan Deck {
	out := make(chan Deck)
	go func() {
		defer close(out)
		for dir := range in {
			if deck, ok := cfg.Deck(dir.Base); ok {
				logger.Debugf("deck filter: forwarding %s with %d files", dir.Base, len(dir.Paths))
				out <- Deck{dir: dir, cfg: deck}
			} else {
				logger.Debugf("deck filter: discarded %s with %d files", dir.Base, len(dir.Paths))
			}
		}
	}()
	return out
}
