package action

import (
	"context"
	"sync"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/worker"
)

// SyncClient interface.
type SyncClient interface {
	worker.Client
	deck.Client
	card.Client
}

// Filesystem interface.
type Filesystem interface {
	worker.Walker
	card.Reader
}

// Parser interface.
type Parser interface {
	card.Parser
	Extensions() []string
}

// Lockfile interface.
type Lockfile interface {
	worker.Lockfile
	deck.Lockfile
	Updated() bool
	Write() error
}

// Sync syncs the cards.
func Sync(ctx context.Context, logger Logger, client SyncClient, fs Filesystem, parser Parser, config deck.Config, lf Lockfile, workspace string) (updated bool, err error) {
	defer func() {
		if writeErr := lf.Write(); err == nil {
			err = writeErr
		}
	}()

	wg := &sync.WaitGroup{}
	errC := make(chan error)
	defer close(errC)

	go func() {
		for err := range errC {
			logger.Errorf("workers: %v", err)
		}
	}()

	dirC, err := worker.FileWalk(ctx, logger, fs, workspace, parser.Extensions())
	if err != nil {
		return false, err
	}

	deckR := worker.SyncDecks(ctx, logger, client, config, lf, dirC)
	deckC := worker.Unwrap(wg, deckR, errC)
	syncR := worker.SyncRequests(ctx, logger, client, fs, parser, lf, workspace, deckC)
	syncC := worker.Unwrap(wg, syncR, errC)
	doneR := worker.ExecuteRequests(ctx, logger, client, lf, syncC)
	_ = worker.Unwrap(wg, doneR, errC)

	wg.Wait()

	return lf.Updated(), err
}
