package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/cmd/github-action/github"
	"github.com/leonhfr/mochi/internal/action"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/parser"
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

	gha.Infof("workspace: %s", workspace)

	fs := file.NewSystem()
	parser := parser.New()
	config, err := action.LoadConfig(fs, gha, parser.List(), workspace)
	if err != nil {
		return err
	}

	client := action.LoadClient(gha, config.RateLimit, token)

	lf, err := action.LoadLockfile(ctx, gha, client, fs, workspace)
	if err != nil {
		return err
	}

	updated, err := action.Sync(ctx, gha, client, fs, parser, config, lf, workspace)
	github.SetOutput(gha, updated)
	return err
}
