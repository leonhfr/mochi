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
	parsers := []parser.Parser{parser.NewNote()}

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

	return Output{}, err
}

type Client interface {
	ListDecks(ctx context.Context) ([]api.Deck, error)
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
