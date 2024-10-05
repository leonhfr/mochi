package main

import (
	"path/filepath"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/internal/action"
)

const (
	apiTokenInput  = "api_token"
	workspaceInput = "workspace"
)

func main() {
	gha := githubactions.New()

	ghc, err := gha.Context()
	if err != nil {
		gha.Fatalf("context error: %v", err)
	}

	token := gha.GetInput(apiTokenInput)
	workspace := gha.GetInput(workspaceInput)

	if err = action.Run(token, filepath.Join(ghc.Workspace, workspace), gha); err != nil {
		gha.Fatalf("error: %v", err)
	}
}
