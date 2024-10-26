package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/heap"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/mochi"
)

// Deck contains a synced deck.
type Deck struct {
	deckID string
	cards  []deck.Card
}

// SyncDecks syncs the decks and parses the files.
func SyncDecks(ctx context.Context, logger Logger, r deck.Reader, p deck.Parser, client deck.CreateClient, config deck.CreateConfig, lf deck.CreateLockfile, workspace string, in <-chan heap.Group[heap.Path]) <-chan Result[Deck] {
	out := make(chan Result[Deck])
	go func() {
		defer close(out)
		for group := range in {
			deckConfig, ok := config.Deck(group.Base)
			if !ok {
				logger.Infof("parse(%s): discarding deck", group.Base)
				continue
			}

			logger.Infof("parse(%s): creating missing decks", group.Base)
			deckID, err := deck.Create(ctx, client, config, lf, group.Base)
			if err != nil {
				out <- Result[Deck]{err: err}
				continue
			}

			logger.Infof("parse(%s): parsing %d files", group.Base, len(group.Items))
			filePaths := heap.ConvertPaths(group.Items)
			cards, err := deck.Parse(r, p, workspace, deckConfig.Parser, filePaths)
			if err != nil {
				out <- Result[Deck]{err: err}
				continue
			}

			deckHeap := deck.Heap(cards)
			logger.Infof("parse(%s): parsed %d cards into %s decks", group.Base, len(cards), deckHeap.Len())
			for deckHeap.Len() > 0 {
				group := deckHeap.Pop()
				out <- Result[Deck]{
					data: Deck{
						deckID: deckID,
						cards:  group.Items,
					},
				}
			}
		}
	}()
	return out
}

// Client is the interface the mochi client should implement to generate the sync requests.
type Client interface {
	ListCardsInDeck(ctx context.Context, deckID string) ([]mochi.Card, error)
}

// SyncRequests returns a stream of requests to sync the cards.
func SyncRequests(ctx context.Context, logger Logger, client Client, lf deck.SyncLockfile, in <-chan Deck) <-chan Result[request.Request] {
	out := make(chan Result[request.Request], inflightRequests)
	go func() {
		defer close(out)

		s := stream.New()
		for deck := range in {
			deck := deck
			s.Go(func() stream.Callback {
				reqs, err := syncRequests(ctx, logger, client, lf, deck)
				if err != nil {
					return func() { out <- Result[request.Request]{err: err} }
				}

				return func() {
					for _, req := range reqs {
						out <- Result[request.Request]{data: req}
					}
				}
			})
		}
		s.Wait()
	}()

	return out
}

func syncRequests(ctx context.Context, logger Logger, client Client, lf deck.SyncLockfile, syncDeck Deck) ([]request.Request, error) {
	logger.Infof("card sync(deckID %s): fetching cards", syncDeck.deckID)
	mochiCards, err := client.ListCardsInDeck(ctx, syncDeck.deckID)
	if err != nil {
		return nil, err
	}
	logger.Infof("card sync(deckID %s): %d cards found", syncDeck.deckID, len(mochiCards))

	cards := deck.ConvertCards(syncDeck.cards)
	logger.Infof("card sync(deckID %s): generating sync requests", syncDeck.deckID)
	reqs := deck.SyncRequests(lf, syncDeck.deckID, mochiCards, cards)
	return reqs, nil
}
