package action

import (
	"context"
	"errors"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
)

const apiTokenInput = "api_token"

type Input struct {
	APIToken  string
	Workspace string
}

func GetInput(gha *githubactions.Action) (Input, error) {
	apiToken := gha.GetInput(apiTokenInput)
	if apiToken == "" {
		return Input{}, errors.New("api token required")
	}

	ghc, err := gha.Context()
	if err != nil {
		return Input{}, err
	}

	return Input{
		APIToken:  apiToken,
		Workspace: ghc.Workspace,
	}, nil
}

func Run(_ context.Context, _ *githubactions.Action, _ Client, _ filesystem.Interface) error {
	return nil
}

type Client interface {
	ListDecks(ctx context.Context) ([]api.Deck, error)
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
