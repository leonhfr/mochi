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
	"github.com/leonhfr/mochi/test/data/base64"
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
			cards     map[string][]api.Card
		}

		want struct {
			deckCreates map[api.CreateDeckRequest]api.Deck
			deckUpdates map[string]api.UpdateDeckRequest
			cardCreates map[string]api.CreateCardRequest
			cardUpdates map[string]api.UpdateCardRequest
			lockFile    string
			output      Output
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
				[]api.Template{
					{
						ID: "xxxxxxxx",
						Fields: map[string]api.FieldTemplate{
							"aaaaaaaa": {ID: "aaaaaaaa"},
							"bbbbbbbb": {ID: "bbbbbbbb"},
							"cccccccc": {ID: "cccccccc"},
						},
					},
				},
				[]api.Deck{
					{
						Name: "Notes (root)",
						ID:   "id_root",
					},
				},
				map[string][]api.Card{
					"id_root": {
						{
							ID:      "id_card_note_2",
							DeckID:  "id_root",
							Name:    "Note 2",
							Content: "# Note 2\n",
						},
					},
					"id_german_vocabulary": {},
				},
			},
			want{
				deckCreates: map[api.CreateDeckRequest]api.Deck{
					{Name: "German"}: {ID: "id_german"},
					{Name: "Vocabulary", ParentID: "id_german"}: {ID: "id_german_vocabulary"},
				},
				deckUpdates: map[string]api.UpdateDeckRequest{},
				cardCreates: map[string]api.CreateCardRequest{
					"id_card_vocabulary_1": {
						DeckID:     "id_german_vocabulary",
						TemplateID: "xxxxxxxx",
						Fields: map[string]api.Field{
							"aaaaaaaa": {ID: "aaaaaaaa", Value: "Spaziergang"},
							"bbbbbbbb": {ID: "bbbbbbbb"},
							"cccccccc": {ID: "cccccccc"},
						},
					},
					"id_card_note_1": {
						Content: "# Note 1\n\nA simple note.\n\n![Constellations](@media/475e4a2888d507e5.jpg)\n\n![Scream](@media/d9bc5d59efbd3aca.png)\n",
						DeckID:  "id_root",
						Attachments: []api.Attachment{
							{
								FileName:    "475e4a2888d507e5.jpg",
								ContentType: "image/jpg",
								Data:        string(base64.Constellations),
							},
						},
					},
				},
				cardUpdates: map[string]api.UpdateCardRequest{
					"id_card_note_2": {
						Content: "# Note 2\n\nAdding an image to an existing note.\n\n![Scream](@media/d9bc5d59efbd3aca.png)\n",
						DeckID:  "id_root",
						Attachments: []api.Attachment{
							{
								FileName:    "d9bc5d59efbd3aca.png",
								ContentType: "image/png",
								Data:        string(base64.Scream),
							},
						},
					},
					"id_card_note_1": {
						Content: "# Note 1\n\nA simple note.\n\n![Constellations](@media/475e4a2888d507e5.jpg)\n\n![Scream](@media/d9bc5d59efbd3aca.png)\n",
						DeckID:  "id_root",
						Attachments: []api.Attachment{
							{
								FileName:    "d9bc5d59efbd3aca.png",
								ContentType: "image/png",
								Data:        string(base64.Scream),
							},
						},
					},
				},
				lockFile: `{"id_german":{"path":"/german","name":""},"id_german_vocabulary":{"path":"/german/vocabulary","name":"","cards":{"id_card_vocabulary_1":{"filename":"s.md"}}},"id_root":{"path":"/","name":"Notes (root)","cards":{"id_card_note_1":{"filename":"note-1.md","images":{"/images/constellations.jpg":"d76091ac3aa97c6fa44f05e35f848332","/images/scream.png":"637b04d6cbd2a4a365fe57c16c90a046"}},"id_card_note_2":{"filename":"note-2.md","images":{"/images/scream.png":"637b04d6cbd2a4a365fe57c16c90a046"}}}}}`,
				output:   Output{LockFileUpdated: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Client
			client := new(MockClient)
			client.On("ListDecks", mock.Anything).Return(tt.api.decks, nil)
			client.On("ListTemplates", mock.Anything).Return(tt.api.templates, nil)
			for id, cards := range tt.api.cards {
				client.On("ListCardsInDeck", mock.Anything, id).Return(cards, nil)
			}
			for req, deck := range tt.want.deckCreates {
				client.On("CreateDeck", mock.Anything, req).Return(deck, nil)
			}
			for id, req := range tt.want.deckUpdates {
				client.On("UpdateDeck", mock.Anything, id, req).Return(api.Deck{}, nil)
			}
			for id, req := range tt.want.cardCreates {
				client.On("CreateCard", mock.Anything, req).Return(api.Card{ID: id}, nil)
			}
			for id, req := range tt.want.cardUpdates {
				client.On("UpdateCard", mock.Anything, id, req).Return(api.Card{ID: id}, nil)
			}

			// Filesystem
			fs := new(MockFilesystem)
			fs.On("Write", "mochi.lock", tt.want.lockFile).Return(nil)

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

func (m *MockFilesystem) Image(path string) ([]byte, string, error) {
	fs := filesystem.New(workspace)
	return fs.Image(path)
}

func (m *MockFilesystem) Sources(extensions []string) ([]string, error) {
	fs := filesystem.New(workspace)
	return fs.Sources(extensions)
}
