package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// ExecuteRequests executes the sync requests.
func ExecuteRequests(ctx context.Context, logger Logger, client *mochi.Client, lf *lock.Lock, in <-chan card.SyncRequest) <-chan Result[struct{}] {
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
