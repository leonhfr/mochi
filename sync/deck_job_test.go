package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonhfr/mochi/parser"
)

func Test_newJobMap(t *testing.T) {
	type (
		input struct {
			config  Config
			sources []string
			lock    *Lock
		}
		want map[string]*deckJob
	)

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			"valid",
			input{
				Config{
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
					Templates: map[string]Template{
						"german": {
							Parser:     "vocabulary",
							TemplateID: "xxxxxxx",
							Fields: map[string]string{
								"aaaaaaaa": "word",
								"bbbbbbbb": "examples",
								"cccccccc": "notes",
							},
						},
					},
				},
				[]string{
					"/german/vocabulary/s.md",
					"/german/vocabulary/f.md",
					"/german/vocabulary/p.md",
					"/german/grammar/gender.md",
					"/note.md",
				},
				&Lock{
					Decks: map[string][2]string{
						"/":                  {"aaaaaaaa", ""},
						"/german":            {"bbbbbbbb", ""},
						"/german/vocabulary": {"cccccccc", ""},
						"/german/grammar":    {"dddddddd", ""},
					},
				},
			},
			map[string]*deckJob{
				"/": {
					sources: []string{"/note.md"},
					id:      "aaaaaaaa",
					parser:  parser.NewNote(),
				},
				"/german/grammar": {
					sources: []string{"/german/grammar/gender.md"},
					id:      "dddddddd",
					parser:  parser.NewNote(),
				},
				"/german/vocabulary": {
					sources:     []string{"/german/vocabulary/s.md", "/german/vocabulary/f.md", "/german/vocabulary/p.md"},
					id:          "cccccccc",
					parser:      parser.NewVocabulary(),
					hasTemplate: true,
					template: Template{
						Parser:     "vocabulary",
						TemplateID: "xxxxxxx",
						Fields:     map[string]string{"aaaaaaaa": "word", "bbbbbbbb": "examples", "cccccccc": "notes"},
					},
				},
			},
		},
	}

	parsers := []parser.Parser{parser.NewNote(), parser.NewVocabulary()}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobMap, err := newJobMap(parsers, tt.input.sources, tt.input.lock, tt.input.config)

			require.NoError(t, err)
			require.Equal(t, len(tt.want), len(jobMap))

			for path, got := range jobMap {
				want, ok := tt.want[path]
				require.True(t, ok)

				assert.Equal(t, want.sources, got.sources)
				assert.Equal(t, want.id, got.id)
				assert.Equal(t, want.archive, got.archive)
				assert.Equal(t, want.hasTemplate, got.hasTemplate)
				assert.Equal(t, want.sources, got.sources)
				assert.IsType(t, want.parser, got.parser)
				assert.Equal(t, want.template, got.template)
			}
		})
	}
}
