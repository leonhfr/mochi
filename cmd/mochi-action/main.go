package main

import (
	"context"
	"time"

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

	transport := action.NewThrottledTransport(30, 15*time.Second)
	client := api.New(input.APIToken, api.WithTransport(transport))
	fs := filesystem.New(input.Workspace)

	output, err := action.Run(ctx, input.ChangedFiles, gha, client, fs)
	if err != nil {
		gha.Fatalf("%v", err)
	}

	action.SetOutput(gha, output)
}
