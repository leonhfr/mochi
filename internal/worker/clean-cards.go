package worker

import (
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// CleanedCards contains a deck with existing cards.
type CleanedCards struct {
	deckID    string
	filePaths []string
	cards     []mochi.Card
}

// CleanCards cleans the lockfile from the cards that have been removed.
func CleanCards(logger Logger, lf *lock.Lock, in <-chan ExistingCards) <-chan CleanedCards {
	out := make(chan CleanedCards)
	go func() {
		defer close(out)

		for existingCards := range in {
			logger.Infof("cleaning cards for deck %s", existingCards.deckID)
			var cardIDs []string
			for _, card := range existingCards.cards {
				cardIDs = append(cardIDs, card.ID)
			}
			lf.CleanCards(existingCards.deckID, cardIDs)
			out <- CleanedCards(existingCards)
		}
	}()
	return out
}
