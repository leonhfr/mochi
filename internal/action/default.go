package action

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/worker"
)

// Logger is the interface to log output.
type Logger interface {
	Infof(format string, args ...any)
}

// Run runs the default action.
func Run(ctx context.Context, logger Logger, token, workspace string) error {
	if token == "" {
		return fmt.Errorf("api token required")
	}

	fs := file.NewSystem()
	cfg, err := config.Parse(fs, workspace)
	if err != nil {
		return err
	}

	lf, err := lock.Parse(fs, workspace)
	if err != nil {
		return err
	}

	dirc, err := worker.FileWalk(ctx, fs, workspace, []string{".md"})
	if err != nil {
		return err
	}

	dirs := []deck.Directory{}
	for dir := range dirc {
		dirs = append(dirs, dir)
	}

	logger.Infof("Hello, world!")
	logger.Infof("api token: %s", token)
	logger.Infof("workspace: %s", workspace)
	logger.Infof("config: %v", cfg)
	logger.Infof("lockfile: %v", lf.String())
	logger.Infof("dirs: %v", dirs)
	return nil
}
