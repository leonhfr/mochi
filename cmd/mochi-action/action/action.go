package action

import (
	"context"
	"errors"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
	"github.com/leonhfr/mochi/sync"
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

func Run(ctx context.Context, gha *githubactions.Action, client Client, fs filesystem.Interface) error {
	parsers := []parser.Parser{parser.NewNote()}

	gha.Noticef("Reading config...")
	config, err := sync.ReadConfig(ctx, parsers, client, fs)
	if err != nil {
		return err
	}

	gha.Noticef("Reading lock file...")
	lock, err := sync.ReadLock(ctx, client, fs)
	if err != nil {
		return err
	}

	defer func() {
		gha.Noticef("Writing lock file...")
		_, terr := lock.Write(fs)
		if terr != nil {
			err = terr
		}
	}()

	gha.Noticef("Searching for sources...")
	sources, err := sync.Sources(config, fs)
	if err != nil {
		return err
	}
	gha.Noticef("%d sources found!", len(sources))

	return nil
}

type Client interface {
	ListDecks(ctx context.Context) ([]api.Deck, error)
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
