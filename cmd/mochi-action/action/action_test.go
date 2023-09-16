package action

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
)

var _ Client = &api.Client{}

var workspace = "../../../test/data"

func Test_GetInput(t *testing.T) {
	//nolint:gosec
	token, workspace := "mochi_cards_token", "/mochi"
	changedFiles := []string{"/mochi/note-1.md", "/mochi/note-2.md"}

	tests := []struct {
		name   string
		envMap map[string]string
		want   Input
		err    string
	}{
		{
			name: "working",
			envMap: map[string]string{
				"INPUT_API_TOKEN":  token,
				"GITHUB_WORKSPACE": workspace,
			},
			want: Input{
				APIToken:  token,
				Workspace: workspace,
			},
			err: "",
		},
		{
			name: "optional",
			envMap: map[string]string{
				"INPUT_API_TOKEN":     token,
				"INPUT_CHANGED_FILES": strings.Join(changedFiles, changedFilesSeparator),
				"GITHUB_WORKSPACE":    workspace,
			},
			want: Input{
				APIToken:     token,
				Workspace:    workspace,
				ChangedFiles: changedFiles,
			},
			err: "",
		},
		{
			name: "missing api token",
			envMap: map[string]string{
				"GITHUB_WORKSPACE": workspace,
			},
			want: Input{},
			err:  "api token required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string {
				return tt.envMap[key]
			}

			gha := githubactions.New(
				githubactions.WithWriter(io.Discard),
				githubactions.WithGetenv(getenv),
			)

			got, err := GetInput(gha)

			assert.Equal(t, tt.want, got)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
		})
	}
}

func Test_Run(t *testing.T) {
	type (
		apiResponses struct {
			templates []api.Template
			decks     []api.Deck
		}

		want struct {
			lockFile string
			output   Output
		}
	)

	tests := []struct {
		name         string
		changedFiles []string
		api          apiResponses
		want         want
	}{
		{
			"all files",
			nil,
			apiResponses{
				nil,
				[]api.Deck{
					{
						Name: "Notes (root)",
						ID:   "id_root",
					},
				},
			},
			want{
				"[decks]\n\"/\" = [\"id_root\", \"Notes (root)\"]\n",
				Output{LockFileUpdated: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Client
			client := new(MockClient)
			client.On("ListDecks", mock.Anything).Return(tt.api.decks, nil)
			client.On("ListTemplates", mock.Anything).Return(tt.api.templates, nil)

			// Filesystem
			fs := new(MockFilesystem)
			fs.On("Write", "mochi-lock.toml", tt.want.lockFile).Return(nil)

			// Run
			ctx := context.Background()
			gha := githubactions.New(
				githubactions.WithWriter(io.Discard),
			)

			output, err := Run(ctx, tt.changedFiles, gha, client, fs)

			require.NoError(t, err)
			assert.Equal(t, tt.want.output, output)
			client.AssertExpectations(t)
			fs.AssertExpectations(t)
		})
	}
}

type MockClient struct {
	mock.Mock
}

var _ Client = &MockClient{}

func (m *MockClient) CreateCard(ctx context.Context, req api.CreateCardRequest) (api.Card, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(api.Card), args.Error(1)
}

func (m *MockClient) ListCardsInDeck(ctx context.Context, id string) ([]api.Card, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockClient) UpdateCard(ctx context.Context, id string, req api.UpdateCardRequest) (api.Card, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(api.Card), args.Error(1)
}

func (m *MockClient) ListDecks(ctx context.Context) ([]api.Deck, error) {
	args := m.Called(ctx)
	return args.Get(0).([]api.Deck), args.Error(1)
}

func (m *MockClient) CreateDeck(ctx context.Context, req api.CreateDeckRequest) (api.Deck, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(api.Deck), args.Error(1)
}

func (m *MockClient) UpdateDeck(ctx context.Context, id string, req api.UpdateDeckRequest) (api.Deck, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(api.Deck), args.Error(1)
}

func (m *MockClient) ListTemplates(ctx context.Context) ([]api.Template, error) {
	args := m.Called(ctx)
	return args.Get(0).([]api.Template), args.Error(1)
}

var _ filesystem.Interface = &MockFilesystem{}

type MockFilesystem struct {
	mock.Mock
}

func (m *MockFilesystem) FileExists(path string) bool {
	fs := filesystem.New(workspace)
	return fs.FileExists(path)
}

func (m *MockFilesystem) Read(path string) ([]byte, error) {
	fs := filesystem.New(workspace)
	return fs.Read(path)
}

func (m *MockFilesystem) Write(path, content string) error {
	args := m.Called(path, content)
	return args.Error(0)
}

func (m *MockFilesystem) Sources(extensions []string) ([]string, error) {
	fs := filesystem.New(workspace)
	return fs.Sources(extensions)
}
