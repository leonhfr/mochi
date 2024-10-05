package action

import (
	"fmt"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/file"
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

	cfg, err := config.Parse(workspace)
	if err != nil {
		return err
	}

	files, err := file.List(workspace, []string{".md"})
	if err != nil {
		return err
	}

	logger.Infof("Hello, world!")
	logger.Infof("api token: %s", token)
	logger.Infof("workspace: %s", workspace)
	logger.Infof("config: %v", cfg)
	logger.Infof("files: %v", files)
	return nil
}
