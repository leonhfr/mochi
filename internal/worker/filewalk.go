package worker

import (
	"container/heap"
	"context"

	"github.com/leonhfr/mochi/internal/deck"
)

// Walker is the interface that should be implemented to recursively walk directories.
type Walker interface {
	Walk(string, []string, func(string)) error
}

// FileWalk is the worker that recursively walks directories and outputs them by
// priority (shorter base directory length).
func FileWalk(ctx context.Context, logger Logger, walker Walker, workspace string, extensions []string) (<-chan deck.Directory, error) {
	h := &deck.Heap{}
	heap.Init(h)

	if err := walker.Walk(
		workspace,
		extensions,
		func(path string) { heap.Push(h, path) },
	); err != nil {
		out := make(chan deck.Directory)
		close(out)
		return out, err
	}

	logger.Infof("filewalk: found %d directories", h.Len())

	out := make(chan deck.Directory)
	go func() {
		defer close(out)
		for h.Len() > 0 {
			select {
			case out <- heap.Pop(h).(deck.Directory):
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}
