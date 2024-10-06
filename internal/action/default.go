package action

import (
	"container/heap"
	"fmt"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
)

// Logger is the interface to log output.
type Logger interface {
	Infof(format string, args ...any)
}

// Run runs the default action.
func Run(token, workspace string, logger Logger) error {
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

	h := &deck.Heap{}
	heap.Init(h)
	err = fs.List(
		workspace,
		[]string{".md"},
		func(path string) { heap.Push(h, path) },
	)
	if err != nil {
		return err
	}

	logger.Infof("Hello, world!")
	logger.Infof("api token: %s", token)
	logger.Infof("workspace: %s", workspace)
	logger.Infof("config: %v", cfg)
	logger.Infof("lockfile: %v", lf.String())
	logger.Infof("files: %v", *h)
	return nil
}
