package main

import (
	"path/filepath"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/internal/config"
)

const (
	apiTokenInput  = "api_token"
	workspaceInput = "workspace"
)

func main() {
	gha := githubactions.New()

	ghc, err := gha.Context()
	if err != nil {
		gha.Fatalf("context error %v", err)
	}

	token := gha.GetInput(apiTokenInput)
	if token == "" {
		gha.Fatalf("%s required", apiTokenInput)
	}

	workspace := gha.GetInput(workspaceInput)

	cfg, err := config.Parse(filepath.Join(ghc.Workspace, workspace))
	if err != nil {
		gha.Fatalf("config error %v", err)
	}

	gha.Infof("Hello, world!")
	gha.Infof("%s=%s", apiTokenInput, token)
	gha.Infof("%s=%s", workspaceInput, workspace)
	gha.Infof("config: %v", cfg)
}
