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
func DeckFilter(logger Logger, cfg *config.Config, dirc <-chan deck.Directory) <-chan Deck {
	deckc := make(chan Deck)
	go func() {
		defer close(deckc)
		for dir := range dirc {
			if deck, ok := cfg.Deck(dir.Base); ok {
				logger.Debugf("deck filter: forwarding %s with %d files", dir.Base, len(dir.Paths))
				deckc <- Deck{dir: dir, cfg: deck}
			}
		}
	}()
	return deckc
}
