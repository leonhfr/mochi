package main

import (
	"context"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/cmd/mochi-action/action"
	"github.com/leonhfr/mochi/filesystem"
)

func main() {
	ctx := context.Background()
	gha := githubactions.New()

	input, err := action.GetInput(gha)
	if err != nil {
		gha.Fatalf("%v", err)
	}

	client := api.New(input.APIToken)
	fs := filesystem.New(input.Workspace)

	if err := action.Run(ctx, gha, client, fs); err != nil {
		gha.Fatalf("%v", err)
	}
}
