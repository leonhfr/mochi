package worker

import (
	"context"

	"github.com/leonhfr/mochi/internal/sync"
)

// Logger is the interface to log output.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

// Walker is the interface that should be implemented to recursively walk directories.
type Walker interface {
	Walk(string, []string, func(string)) error
}

// FileWalk is the worker that recursively walks directories and outputs them by
// priority (shorter base directory length).
func FileWalk(ctx context.Context, logger Logger, walker Walker, workspace string, extensions []string) (<-chan sync.Directory, error) {
	h := sync.NewHeap()

	if err := walker.Walk(
		workspace,
		extensions,
		func(path string) { h.Push(path) },
	); err != nil {
		out := make(chan sync.Directory)
		close(out)
		return out, err
	}

	logger.Infof("filewalk: found %d directories", h.Len())

	out := make(chan sync.Directory, h.Len())
	go func() {
		defer close(out)

		for h.Len() > 0 {
			select {
			case out <- h.Pop():
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}
