package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/mochi"
)

const inflightRequests = 50

// Client is the interface the mochi client should implement to generate the sync requests.
type Client interface {
	ListCardsInDeck(ctx context.Context, deckID string) ([]mochi.Card, error)
}

// Lockfile is the interface the lockfile should implement to generate the sync requests.
type Lockfile interface {
	deck.SyncLockfile
	Deck(id string) (lock.Deck, bool)
	DeleteCard(deckID, cardID string)
}

// SyncRequests returns a stream of requests to sync the cards.
func SyncRequests(ctx context.Context, logger Logger, client Client, lf Lockfile, in <-chan Deck) <-chan Result[request.Request] {
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

// DumpRequests returns a stream of requests to delete the cards.
func DumpRequests(ctx context.Context, logger Logger, client deck.DumpClient, in <-chan string) <-chan Result[request.Request] {
	out := make(chan Result[request.Request], inflightRequests)
	go func() {
		defer close(out)

		s := stream.New()
		for deckID := range in {
			deckID := deckID
			s.Go(func() stream.Callback {
				logger.Infof("dump(deckID %s): generating delete requests", deckID)
				reqs, err := deck.DumpRequests(ctx, client, deckID)
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

// ExecuteRequests executes the sync requests.
func ExecuteRequests(ctx context.Context, logger Logger, client request.Client, reader image.Reader, lf request.Lockfile, in <-chan request.Request) <-chan Result[struct{}] {
	out := make(chan Result[struct{}])
	go func() {
		defer close(out)

		s := stream.New()
		for req := range in {
			req := req
			s.Go(func() stream.Callback {
				logger.Infof("syncing: %s", req.String())
				if err := req.Execute(ctx, client, reader, lf); err != nil {
					return func() {
						out <- Result[struct{}]{err: err}
					}
				}
				return func() {}
			})
		}
		s.Wait()
	}()

	return out
}

func syncRequests(ctx context.Context, logger Logger, client Client, lf Lockfile, syncDeck Deck) ([]request.Request, error) {
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
