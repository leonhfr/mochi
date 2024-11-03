package worker

import (
	"context"

	"github.com/sourcegraph/conc/stream"

	"github.com/leonhfr/mochi/internal/request"
)

const inflightRequests = 50

// ExecuteRequests executes the sync requests.
func ExecuteRequests(ctx context.Context, logger Logger, client request.Client, lf request.Lockfile, in <-chan request.Request) <-chan Result[struct{}] {
	out := make(chan Result[struct{}])
	go func() {
		defer close(out)

		s := stream.New()
		for req := range in {
			req := req
			s.Go(func() stream.Callback {
				logger.Infof("executing: %s", req.String())
				if err := req.Execute(ctx, client, lf); err != nil {
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
