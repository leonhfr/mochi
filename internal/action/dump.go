package action

import (
	"context"
	"sync"

	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/internal/worker"
)

// Dump deletes all the cards and decks.
func Dump(ctx context.Context, logger Logger, token, workspace string) (err error) {
	logger.Infof("workspace: %s", workspace)

	fs := file.NewSystem()
	config, err := loadConfig(fs, logger, parser.Names(), workspace)
	if err != nil {
		return err
	}

	client := loadClient(logger, config.RateLimit, token)

	wg := &sync.WaitGroup{}
	errC := make(chan error)
	defer close(errC)

	go func() {
		for err := range errC {
			logger.Errorf("workers: %v", err)
		}
	}()

	deckC, err := worker.ListDecks(ctx, client)
	if err != nil {
		return err
	}

	lf := &noOpLockfile{}

	dumpR := worker.DumpRequests(ctx, logger, client, deckC)
	dumpC := worker.Unwrap(wg, dumpR, errC)
	doneR := worker.ExecuteRequests(ctx, logger, client, lf, dumpC)
	_ = worker.Unwrap(wg, doneR, errC)

	wg.Wait()

	for {
		ok, err := worker.DeleteLeafDecks(ctx, client)
		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	return err
}

var _ request.Lockfile = &noOpLockfile{}

type noOpLockfile struct{}

func (lf *noOpLockfile) Lock()   {}
func (lf *noOpLockfile) Unlock() {}

func (lf *noOpLockfile) SetCard(_, _, _ string) error {
	return nil
}
