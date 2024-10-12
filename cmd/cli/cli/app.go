package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/leonhfr/mochi/internal/action"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/parser"
)

// GetApp returns the cli app.
func GetApp(out io.Writer, version, compiled string) (*cli.App, error) {
	compiledTime, err := time.Parse(time.RFC3339, compiled)
	if err != nil {
		return nil, err
	}

	logger := newLogger(out)

	return &cli.App{
		Name:      "mochi",
		Usage:     "synchronize markdown notes to mochi cards",
		Version:   version,
		Compiled:  compiledTime,
		Args:      true,
		ArgsUsage: "[workspace]",
		Before: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return fmt.Errorf("expected one argument")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}

			token := ctx.String("token")
			workspace := ctx.Args().First()
			workspace = filepath.Join(pwd, workspace)

			logger.Infof("workspace: %s", workspace)

			fs := file.NewSystem()
			parser := parser.New()
			config, err := action.LoadConfig(fs, logger, parser.List(), workspace)
			if err != nil {
				return err
			}

			client := action.LoadClient(logger, config.RateLimit, token)

			lf, err := action.LoadLockfile(ctx.Context, logger, client, fs, workspace)
			if err != nil {
				return err
			}

			_, err = action.Sync(ctx.Context, logger, client, fs, parser, config, lf, workspace)
			return err
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "mochi API token",
				EnvVars: []string{"MOCHI_API_TOKEN"},
			},
		},
	}, nil
}
