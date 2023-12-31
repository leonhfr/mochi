package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/parser"
)

var config = `sync:
  - path: german/vocabulary
    template: german
    archive: true
  - path: .
    name: Notes (root)
    parser: note
    archive: true
    walk: true
  - path: german/grammar
    parser: note
    archive: true

ignore:
  - journal/*

templates:
  german:
    parser: vocabulary
    templateId: xxxxxxxx
    fields:
      aaaaaaaa: word
      bbbbbbbb: examples
      cccccccc: notes
`

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

	fs := new(MockFilesystem)
	fs.On("FileExists", "mochi.yaml").Return(false)
	fs.On("FileExists", "mochi.yml").Return(true)
	fs.On("Read", "mochi.yml").Return([]byte(config), nil)
	got, err := ReadConfig(context.Background(), parsers, client, fs)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	fs.AssertExpectations(t)
}

func Test_Config_deckName(t *testing.T) {
	tests := []struct {
		sync []Sync
		path string
		want string
	}{
		{[]Sync{{Path: "/", Name: "Notes (root)"}}, "/", "Notes (root)"},
		{[]Sync{}, "/", rootDeckName},
		{[]Sync{}, "/journal", "Journal"},
		{[]Sync{}, "/german/vocabulary", "Vocabulary"},
		{[]Sync{{Path: "German", Name: "Deutsch"}}, "/german/vocabulary", "Vocabulary"},
		{[]Sync{}, "/sub/developer-experience", "Developer Experience"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			config := Config{Sync: tt.sync}
			got := config.deckName(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Config_ignored(t *testing.T) {
	config := Config{Ignore: []string{
		"/journal/**",
	}}

	tests := []struct {
		path string
		want bool
	}{
		{"/journal/yyyy-mm-dd.md", true},
		{"/german/vocabulary/s.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := config.ignored(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Config_matchSync(t *testing.T) {
	config := Config{Sync: []Sync{
		{
			Path:   "/german/vocabulary",
			Parser: "vocabulary",
		},
		{
			Path:   "/german/grammar",
			Parser: "grammar",
		},
		{
			Path:   "/sub",
			Parser: "root",
			Walk:   true,
		},
	}}

	tests := []struct {
		path   string
		want   bool
		parser string // sync parser
	}{
		{"/german/vocabulary/s.md", true, "vocabulary"},
		{"/german/grammar/noun.md", true, "grammar"},
		{"/sub/note.md", true, "note"},
		{"/sub/sub/note.md", true, "note"},
		{"/note.md", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			_, got := config.matchSync(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
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

	fs := new(MockFilesystem)
	fs.On("Read", "mochi.yml").Return([]byte(config), nil)
	got, err := parseConfig(Config{}, "mochi.yml", fs)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	fs.AssertExpectations(t)
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
