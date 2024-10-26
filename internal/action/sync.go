package action

import (
	"context"
	"sync"

	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/worker"
)

// Sync syncs the cards.
func Sync(ctx context.Context, logger Logger, token, workspace string) (updated bool, err error) {
	logger.Infof("workspace: %s", workspace)

	fs := file.NewSystem()
	parser := parser.New()
	config, err := loadConfig(fs, logger, parser.List(), workspace)
	if err != nil {
		return false, err
	}

	client := loadClient(logger, config.RateLimit, token)

	lf, err := loadLockfile(ctx, logger, client, fs, workspace)
	if err != nil {
		return false, err
	}

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

	deckR := worker.SyncDecks(ctx, logger, fs, parser, client, config, lf, workspace, dirC)
	deckC := worker.Unwrap(wg, deckR, errC)
	syncR := worker.SyncRequests(ctx, logger, client, lf, deckC)
	syncC := worker.Unwrap(wg, syncR, errC)
	doneR := worker.ExecuteRequests(ctx, logger, client, fs, lf, syncC)
	_ = worker.Unwrap(wg, doneR, errC)

	wg.Wait()

	return lf.Updated(), err
}
