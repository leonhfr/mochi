package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/request"
)

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
