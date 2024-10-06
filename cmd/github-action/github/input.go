package github

import (
	"fmt"
	"path/filepath"

	"github.com/sethvargo/go-githubactions"
)

const (
	apiTokenInput  = "api_token"
	workspaceInput = "workspace"
)

// GetInput returns the action inputs.
func GetInput(gha *githubactions.Action) (string, string, error) {
	ghc, err := gha.Context()
	if err != nil {
		return "", "", err
	}

	token := gha.GetInput(apiTokenInput)
	if token == "" {
		return "", "", fmt.Errorf("%s required", apiTokenInput)
	}

	workspace := gha.GetInput(workspaceInput)
	workspace = filepath.Join(ghc.Workspace, workspace)
	return token, workspace, nil
}
