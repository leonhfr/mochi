package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
)

var workspace = "../test/data"

var _ Client = &api.Client{}

func Test_ReadConfig(t *testing.T) {
	templates := []api.Template{
		{
			ID: "xxxxxxxx",
			Fields: map[string]api.FieldTemplate{
				"aaaaaaaa": {ID: "aaaaaaaa"},
				"bbbbbbbb": {ID: "bbbbbbbb"},
				"cccccccc": {ID: "cccccccc"},
			},
		},
	}

	parsers := []parser.Parser{parser.NewNote(), parser.NewVocabulary()}

	want := Config{
		Sync: []Sync{
			{
				Path:     "/german/vocabulary",
				Template: "german",
				Archive:  true,
			},
			{
				Path:    "/german/grammar",
				Parser:  "note",
				Archive: true,
			},
			{
				Path:    "/",
				Name:    "Notes (root)",
				Parser:  "note",
				Archive: true,
				Walk:    true,
			},
		},
		Ignore: []string{
			"/journal/*",
		},
		Templates: map[string]Template{
			"german": {
				Parser:     "vocabulary",
				TemplateID: "xxxxxxxx",
				Fields: map[string]string{
					"aaaaaaaa": "word",
					"bbbbbbbb": "examples",
					"cccccccc": "notes",
				},
			},
		},
		parsers: parsers,
	}

	client := new(MockClient)
	client.On("ListTemplates", mock.Anything).Return(templates, nil)

	fs := filesystem.New(workspace)
	got, err := ReadConfig(context.Background(), parsers, client, fs)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_parseConfig(t *testing.T) {
	want := Config{
		Sync: []Sync{
			{
				Path:     "german/vocabulary",
				Template: "german",
				Archive:  true,
			},
			{
				Path:    ".",
				Name:    "Notes (root)",
				Parser:  "note",
				Archive: true,
				Walk:    true,
			},
			{
				Path:    "german/grammar",
				Parser:  "note",
				Archive: true,
			},
		},
		Ignore: []string{
			"journal/*",
		},
		Templates: map[string]Template{
			"german": {
				Parser:     "vocabulary",
				TemplateID: "xxxxxxxx",
				Fields: map[string]string{
					"aaaaaaaa": "word",
					"bbbbbbbb": "examples",
					"cccccccc": "notes",
				},
			},
		},
	}

	fs := filesystem.New(workspace)
	got, err := parseConfig(Config{}, "mochi.yml", fs)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_cleanConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Config
	}{
		{
			"ok",
			Config{
				Sync: []Sync{
					{Path: ""},
					{Path: "."},
					{Path: "/"},
					{Path: "/german"},
					{Path: "/german/vocabulary"},
					{Path: "/german/vocabulary/"},
					{Path: "german"},
					{Path: "german/vocabulary"},
					{Path: "german/vocabulary/"},
					{Path: "../german/vocabulary/"},
				},
				Ignore: []string{
					"journal",
					"journal/*",
					"/journal/*",
				},
			},
			Config{
				Sync: []Sync{
					{Path: "/german/vocabulary"},
					{Path: "/german/vocabulary"},
					{Path: "/german/vocabulary"},
					{Path: "/german/vocabulary"},
					{Path: "/german/vocabulary"},
					{Path: "/german"},
					{Path: "/german"},
					{Path: "/"},
					{Path: "/"},
					{Path: "/"},
				},
				Ignore: []string{
					"/journal",
					"/journal/*",
					"/journal/*",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanConfig(tt.config)
			assert.Equal(t, tt.want, tt.config)
		})
	}
}

type MockClient struct {
	mock.Mock
}

var _ Client = &MockClient{}

func (m *MockClient) ListTemplates(ctx context.Context) ([]api.Template, error) {
	args := m.Called(ctx)
	return args.Get(0).([]api.Template), args.Error(1)
}
