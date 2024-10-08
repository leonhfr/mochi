package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

// ExistingCards contains a deck with existing cards.
type ExistingCards struct {
	deckID    string
	filePaths []string
	cards     []mochi.Card
}

// FetchCards fetches all the existing cards in a deck.
func FetchCards(ctx context.Context, logger Logger, client *mochi.Client, in <-chan SyncedDeck) <-chan Result[ExistingCards] {
	out := make(chan Result[ExistingCards])
	go func() {
		defer close(out)

		for syncedDeck := range in {
			logger.Infof("fetch cards for deck %s", syncedDeck.id)
			var cards []mochi.Card
			err := client.ListCardsInDeck(
				ctx,
				syncedDeck.id,
				func(cc []mochi.Card) error { cards = append(cards, cc...); return nil },
			)
			out <- Result[ExistingCards]{
				data: ExistingCards{
					deckID:    syncedDeck.id,
					filePaths: syncedDeck.filePaths,
					cards:     cards,
				},
				err: err,
			}
		}
	}()
	return out
}

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

// ParsedCards contains a deck with existing and parsed cards.
type ParsedCards struct {
	deckID      string
	filePaths   []string
	mochiCards  []mochi.Card
	parsedCards []parser.Card
}

// ParseCards parses the files and converts them to cards.
func ParseCards(logger Logger, fs *file.System, parser *parser.Parser, workspace string, in <-chan CleanedCards) <-chan Result[ParsedCards] {
	out := make(chan Result[ParsedCards])
	go func() {
		defer close(out)

		for cleanedCards := range in {
			logger.Infof("parsing cards for deck %s", cleanedCards.deckID)
			parsedCards, err := card.Parse(fs, parser, workspace, cleanedCards.filePaths)
			out <- Result[ParsedCards]{
				data: ParsedCards{
					deckID:      cleanedCards.deckID,
					filePaths:   cleanedCards.filePaths,
					mochiCards:  cleanedCards.cards,
					parsedCards: parsedCards,
				},
				err: err,
			}
		}
	}()
	return out
}
