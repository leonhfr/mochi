package action

import (
	"context"
	"sync"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/worker"
	"github.com/leonhfr/mochi/mochi"
)

// Logger is the interface to log output.
type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
	Infof(format string, args ...any)
}

// Run runs the default action.
func Run(ctx context.Context, logger Logger, token, workspace string) (updated bool, err error) {
	logger.Infof("workspace: %s", workspace)

	client := mochi.New(token)
	fs := file.NewSystem()
	parser := parser.New()

	cfg, err := config.Parse(fs, workspace, parser.List())
	if err != nil {
		return false, err
	}
	logger.Infof("loaded config")
	logger.Debugf("config: %v", cfg)

	lf, err := getLockfile(ctx, client, fs, workspace)
	if err != nil {
		return false, err
	}

	logger.Infof("loaded lockfile")
	logger.Debugf("lockfile: %v", lf.String())

	defer func() {
		if writeErr := lf.Write(); err == nil {
			err = writeErr
		}
	}()

	err = runWorkers(ctx, logger, client, fs, parser, cfg, lf, workspace)

	return lf.Updated(), err
}

func runWorkers(ctx context.Context, logger Logger, client *mochi.Client, fs *file.System, parser *parser.Parser, cfg *config.Config, lf *lock.Lock, workspace string) error {
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
		return err
	}

	deckR := worker.SyncDecks(ctx, logger, client, cfg, lf, dirC)
	deckC := worker.Unwrap(wg, deckR, errC)
	syncR := worker.CardRequests(ctx, logger, client, fs, parser, lf, workspace, deckC)
	syncC := worker.Unwrap(wg, syncR, errC)
	doneC := worker.ExecuteRequests(ctx, logger, client, lf, syncC)
	end := worker.Unwrap(wg, doneC, errC)

	<-end
	wg.Wait()
	return nil
}

func getLockfile(ctx context.Context, client *mochi.Client, fs *file.System, workspace string) (*lock.Lock, error) {
	lf, err := lock.Parse(fs, workspace)
	if err != nil {
		return nil, err
	}

	var decks []mochi.Deck
	if err := client.ListDecks(ctx, func(dd []mochi.Deck) error {
		decks = append(decks, dd...)
		return nil
	}); err != nil {
		return nil, err
	}

	lf.CleanDecks(decks)
	return lf, nil
}
