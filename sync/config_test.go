package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
)

var workspace = "../test/data"

func Test_ReadConfig(t *testing.T) {
	parsers := []parser.Parser{parser.NewNote(), parser.NewVocabulary()}

	want := Config{
		Sync: []Sync{
			{
				Path:     "/german/vocabulary",
				Template: "german",
			},
			{
				Path:   "/german/grammar",
				Parser: "note",
			},
			{
				Path:   "/",
				Name:   "Notes (root)",
				Parser: "note",
				Walk:   true,
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

	fs := filesystem.New(workspace)
	got, err := ReadConfig(parsers, fs)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func Test_parseConfig(t *testing.T) {
	want := Config{
		Sync: []Sync{
			{
				Path:     "german/vocabulary",
				Template: "german",
			},
			{
				Path:   ".",
				Name:   "Notes (root)",
				Parser: "note",
				Walk:   true,
			},
			{
				Path:   "german/grammar",
				Parser: "note",
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
