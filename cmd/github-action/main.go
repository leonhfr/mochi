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

	return action.Run(ctx, gha, token, workspace)
}
