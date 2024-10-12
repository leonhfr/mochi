package cli

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/leonhfr/mochi/internal/action"
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
		Commands: []*cli.Command{
			{
				Name:  "dump",
				Usage: "deletes all cards and decks",
				Action: func(ctx *cli.Context) error {
					pwd, err := os.Getwd()
					if err != nil {
						return err
					}

					token := ctx.String("token")
					workspace := ctx.Args().First()
					workspace = filepath.Join(pwd, workspace)

					return action.Dump(ctx.Context, logger, token, workspace)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "token",
						Aliases: []string{"t"},
						Usage:   "mochi API token",
						EnvVars: []string{"MOCHI_API_TOKEN"},
					},
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}

			token := ctx.String("token")
			workspace := ctx.Args().First()
			workspace = filepath.Join(pwd, workspace)

			_, err = action.Sync(ctx.Context, logger, token, workspace)
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
