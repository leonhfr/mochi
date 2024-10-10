package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/cmd/github-action/github"
	"github.com/leonhfr/mochi/internal/action"
)

func main() {
	gha := githubactions.New()
	if err := run(context.Background(), gha); err != nil {
		gha.Fatalf("error: %v", err)
	}
}

func run(ctx context.Context, gha *githubactions.Action) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	token, workspace, err := github.GetInput(gha)
	if err != nil {
		return err
	}

	client, fs, parser, config, lf, err := action.Load(ctx, gha, token, workspace)
	if err != nil {
		return err
	}

	updated, err := action.Sync(ctx, gha, client, fs, parser, config, lf, workspace)
	github.SetOutput(gha, updated)
	return err
}
