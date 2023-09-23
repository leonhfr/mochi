package action

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sethvargo/go-githubactions"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
	"github.com/leonhfr/mochi/sync"
)

const (
	apiTokenInput         = "api_token"
	changedFilesInput     = "changed_files"
	changedFilesSeparator = " "
	lockFileUpdatedOutput = "lock_file_updated"
)

type Input struct {
	APIToken     string
	Workspace    string
	ChangedFiles []string
}

type Output struct {
	LockFileUpdated bool
}

func GetInput(gha *githubactions.Action) (Input, error) {
	apiToken := gha.GetInput(apiTokenInput)
	if apiToken == "" {
		return Input{}, errors.New("api token required")
	}

	var changedFiles []string
	if input := gha.GetInput(changedFilesInput); len(input) > 0 {
		changedFiles = strings.Split(input, changedFilesSeparator)
	}

	ghc, err := gha.Context()
	if err != nil {
		return Input{}, err
	}

	return Input{
		APIToken:     apiToken,
		Workspace:    ghc.Workspace,
		ChangedFiles: changedFiles,
	}, nil
}

func SetOutput(gha *githubactions.Action, output Output) {
	gha.SetOutput(lockFileUpdatedOutput, fmt.Sprint(output.LockFileUpdated))
}

func Run(ctx context.Context, changedFiles []string, gha *githubactions.Action, client Client, fs filesystem.Interface) (output Output, err error) {
	parsers := []parser.Parser{
		parser.NewNote(),
		parser.NewVocabulary(),
		parser.NewHeadings(),
	}

	gha.Noticef("Reading config...")
	config, err := sync.ReadConfig(ctx, parsers, client, fs)
	if err != nil {
		return Output{}, err
	}

	gha.Noticef("Reading lock file...")
	lock, err := sync.ReadLock(ctx, client, fs)
	if err != nil {
		return Output{}, err
	}

	defer func() {
		gha.Noticef("Writing lock file...")
		updated, terr := lock.Write(fs)
		output.LockFileUpdated = updated
		if terr != nil {
			err = terr
		}
	}()

	gha.Noticef("Searching for sources...")
	sources, err := sync.Sources(changedFiles, config, fs)
	if err != nil {
		return Output{}, err
	}
	gha.Noticef("%d sources found!", len(sources))

	gha.Noticef("Synchronizing decks...")
	dr, err := sync.SynchronizeDecks(ctx, sources, lock, config, client)
	if err != nil {
		return Output{}, err
	}
	gha.Noticef("Created %d and updated %d decks", dr.Created, dr.Updated)

	gha.Noticef("Synchronizing cards...")
	cr, err := sync.SynchronizeCards(ctx, parsers, sources, lock, config, client, fs)
	if err != nil {
		return Output{}, err
	}
	gha.Noticef("Created %d, updated %d, and archived %d cards", cr.Created, cr.Updated, cr.Archived)

	return Output{}, err
}

type Client interface {
	ListCardsInDeck(ctx context.Context, id string) ([]api.Card, error)
	CreateCard(ctx context.Context, req api.CreateCardRequest) (api.Card, error)
	UpdateCard(ctx context.Context, id string, req api.UpdateCardRequest) (api.Card, error)
	ListDecks(ctx context.Context) ([]api.Deck, error)
	CreateDeck(ctx context.Context, req api.CreateDeckRequest) (api.Deck, error)
	UpdateDeck(ctx context.Context, id string, req api.UpdateDeckRequest) (api.Deck, error)
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
