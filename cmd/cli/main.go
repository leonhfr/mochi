package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/leonhfr/mochi/cmd/cli/cli"
)

var (
	version = "0.0.0"
	date    = "2006-01-02T15:04:05+07:00"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	app, err := cli.GetApp(os.Stdout, version, date)
	if err != nil {
		return err
	}

	return app.RunContext(ctx, os.Args)
}
