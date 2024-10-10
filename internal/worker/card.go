package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

const inflightSyncRequests = 20

// CardRequests returns a stream of requests to sync the cards.
func CardRequests(ctx context.Context, logger Logger, client *mochi.Client, fs *file.System, parser *parser.Parser, lf *lock.Lock, workspace string, in <-chan Deck) <-chan Result[card.SyncRequest] {
	out := make(chan Result[card.SyncRequest], inflightSyncRequests)
	go func() {
		defer close(out)

		s := stream.New().WithMaxGoroutines(cap(in))
		for deck := range in {
			deck := deck
			s.Go(func() stream.Callback {
				reqs, err := syncRequests(ctx, logger, client, fs, parser, lf, workspace, deck)
				if err != nil {
					return func() { out <- Result[card.SyncRequest]{err: err} }
				}

				return func() {
					for _, req := range reqs {
						out <- Result[card.SyncRequest]{data: req}
					}
				}
			})
		}
		s.Wait()
	}()

	return out
}

func syncRequests(ctx context.Context, logger Logger, client *mochi.Client, fs *file.System, parser *parser.Parser, lf *lock.Lock, workspace string, deck Deck) ([]card.SyncRequest, error) {
	logger.Infof("card sync(deckID %s): fetching cards", deck.deckID)
	mochiCards, err := fetchCardsInDeck(ctx, client, deck.deckID)
	if err != nil {
		return nil, err
	}
	logger.Infof("card sync(deckID %s): %s cards found", deck.deckID, len(mochiCards))

	logger.Infof("card sync(deckID %s): cleaning lockfile", deck.deckID)
	cleanCards(lf, deck.deckID, mochiCards)

	logger.Infof("card sync(deckID %s): parsing cards", deck.deckID)
	parsedCards, err := card.Parse(fs, parser, workspace, deck.config.Parser, deck.filePaths)
	if err != nil {
		return nil, err
	}

	logger.Infof("card sync(deckID %s): generating sync requests", deck.deckID)
	reqs := card.SyncRequests(lf, deck.deckID, mochiCards, parsedCards)
	return reqs, nil
}

func fetchCardsInDeck(ctx context.Context, client *mochi.Client, deckID string) ([]mochi.Card, error) {
	var cards []mochi.Card
	err := client.ListCardsInDeck(
		ctx,
		deckID,
		func(cc []mochi.Card) error { cards = append(cards, cc...); return nil },
	)
	return cards, err
}

func cleanCards(lf *lock.Lock, deckID string, mochiCards []mochi.Card) {
	cardIDs := make([]string, 0, len(mochiCards))
	for _, card := range mochiCards {
		cardIDs = append(cardIDs, card.ID)
	}
	lf.CleanCards(deckID, cardIDs)
}
