package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// SyncRequests generates the requests required to sync the deck.
func SyncRequests(logger Logger, lf *lock.Lock, in <-chan ParsedCards) <-chan card.SyncRequest {
	out := make(chan card.SyncRequest)
	go func() {
		defer close(out)

		for parsedCards := range in {
			logger.Infof("generating sync requests for deck %s", parsedCards.deckID)
			reqs := card.SyncRequests(lf, parsedCards.deckID, parsedCards.mochiCards, parsedCards.parsedCards)
			for _, req := range reqs {
				out <- req
			}
		}
	}()
	return out
}

// ExecuteRequests executes the sync requests.
func ExecuteRequests(ctx context.Context, logger Logger, client *mochi.Client, lf *lock.Lock, in <-chan card.SyncRequest) <-chan Result[struct{}] {
	out := make(chan Result[struct{}])
	go func() {
		defer close(out)

		for req := range in {
			logger.Infof("executing %s", req.String())
			if err := req.Sync(ctx, client, lf); err != nil {
				out <- Result[struct{}]{err: err}
			}
		}
	}()
	return out
}
