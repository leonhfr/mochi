package action

import (
	"context"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/worker"
	"github.com/leonhfr/mochi/mochi"
)

// Logger is the interface to log output.
type Logger interface {
	Infof(format string, args ...any)
}

// Run runs the default action.
func Run(ctx context.Context, logger Logger, token, workspace string) (updated bool, err error) {
	client := mochi.New(token)
	fs := file.NewSystem()

	cfg, err := config.Parse(fs, workspace)
	if err != nil {
		return false, err
	}

	lf, err := getLockfile(ctx, client, fs, workspace)
	if err != nil {
		return false, err
	}

	defer func() {
		if writeErr := lf.Write(); err == nil {
			err = writeErr
		}
	}()

	err = runWorkers(ctx, logger, fs, cfg, lf, workspace)

	return lf.Updated(), err
}

func runWorkers(ctx context.Context, logger Logger, fs *file.System, cfg *config.Config, lf *lock.Lock, workspace string) error {
	dirc, err := worker.FileWalk(ctx, fs, workspace, []string{".md"})
	if err != nil {
		return err
	}

	deckc := worker.DeckFilter(cfg, dirc)
	decks := []worker.Deck{}
	for deck := range deckc {
		decks = append(decks, deck)
	}

	logger.Infof("workspace: %s", workspace)
	logger.Infof("config: %v", cfg)
	logger.Infof("lockfile: %v", lf.String())
	logger.Infof("decks: %v", decks)
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
