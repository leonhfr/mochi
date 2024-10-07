package worker

import (
	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
)

// FilteredDeck represents a deck whose config has been matched.
type FilteredDeck struct {
	path      string
	filePaths []string
	name      *string
}

// DeckFilter filters the directories, only forwarding them
// if a deck config has been found.
func DeckFilter(logger Logger, cfg *config.Config, in <-chan deck.Directory) <-chan FilteredDeck {
	out := make(chan FilteredDeck)
	go func() {
		defer close(out)

		for dir := range in {
			if deckConfig, ok := cfg.GetDeck(dir.Path); ok {
				logger.Infof("deck filter: forwarding %s with %d files", dir.Path, len(dir.FilePaths))
				out <- FilteredDeck{
					path:      deckConfig.Path,
					filePaths: dir.FilePaths,
					name:      deckConfig.Name,
				}
			} else {
				logger.Infof("deck filter: discarded %s with %d files", dir.Path, len(dir.FilePaths))
			}
		}
	}()
	return out
}
