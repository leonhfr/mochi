package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/card"
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
	card.Lockfile
	CleanCards(deckID string, cardIDs []string)
}

// SyncRequests returns a stream of requests to sync the cards.
func SyncRequests(ctx context.Context, logger Logger, client Client, cr card.Reader, parser card.Parser, lf Lockfile, workspace string, in <-chan Deck) <-chan Result[request.Request] {
	out := make(chan Result[request.Request], inflightRequests)
	go func() {
		defer close(out)

		s := stream.New().WithMaxGoroutines(cap(in))
		for deck := range in {
			deck := deck
			s.Go(func() stream.Callback {
				reqs, err := syncRequests(ctx, logger, client, cr, parser, lf, workspace, deck)
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
func DumpRequests(ctx context.Context, logger Logger, client card.DumpClient, in <-chan string) <-chan Result[request.Request] {
	out := make(chan Result[request.Request], inflightRequests)
	go func() {
		defer close(out)

		s := stream.New().WithMaxGoroutines(cap(in))
		for deckID := range in {
			deckID := deckID
			s.Go(func() stream.Callback {
				logger.Infof("dump(deckID %s): generating delete requests", deckID)
				reqs, err := card.DumpRequests(ctx, client, deckID)
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
func ExecuteRequests(ctx context.Context, logger Logger, client request.Client, lf request.Lockfile, in <-chan request.Request) <-chan Result[struct{}] {
	out := make(chan Result[struct{}])
	go func() {
		defer close(out)

		s := stream.New().WithMaxGoroutines(cap(in))
		for req := range in {
			req := req
			s.Go(func() stream.Callback {
				logger.Infof("syncing: %s", req.String())
				if err := req.Sync(ctx, client, lf); err != nil {
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

func syncRequests(ctx context.Context, logger Logger, client Client, cr card.Reader, parser card.Parser, lf Lockfile, workspace string, deck Deck) ([]request.Request, error) {
	logger.Infof("card sync(deckID %s): fetching cards", deck.deckID)
	mochiCards, err := client.ListCardsInDeck(ctx, deck.deckID)
	if err != nil {
		return nil, err
	}
	logger.Infof("card sync(deckID %s): %s cards found", deck.deckID, len(mochiCards))

	logger.Infof("card sync(deckID %s): cleaning lockfile", deck.deckID)
	cleanCards(lf, deck.deckID, mochiCards)

	logger.Infof("card sync(deckID %s): parsing cards", deck.deckID)
	parsedCards, err := card.Parse(cr, parser, workspace, deck.config.Parser, deck.filePaths)
	if err != nil {
		return nil, err
	}

	logger.Infof("card sync(deckID %s): generating sync requests", deck.deckID)
	reqs := card.SyncRequests(lf, deck.deckID, mochiCards, parsedCards)
	return reqs, nil
}

func cleanCards(lf Lockfile, deckID string, mochiCards []mochi.Card) {
	cardIDs := make([]string, 0, len(mochiCards))
	for _, card := range mochiCards {
		cardIDs = append(cardIDs, card.ID)
	}
	lf.CleanCards(deckID, cardIDs)
}
