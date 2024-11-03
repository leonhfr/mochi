package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/config"
)

var vocabularySource = `<!-- Generated. -->

Spaziergang
"Wir haben nach dem Essen einen langen Spaziergang gemacht."
Stroll, walk, promenade.

Spiegel

First line can be a sentence.

Word
"Example without note for word."

AnotherWord
Notes without example.
Notes can be multiline.`

func Test_vocabulary_parse(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		config config.VocabularyTemplate
		source string
		want   Result
	}{
		{
			name:   "should parse vocabulary",
			path:   "/testdata/languages/de/vocabulary/s.md",
			config: config.VocabularyTemplate{TemplateID: "GERMAN_TEMPLATE"},
			source: "<!-- Generated. -->\n\nSpaziergang\n\nSpiegel\n",
			want: Result{Cards: []Card{
				{
					Fields: map[string]string{
						"name": "Spaziergang",
					},
					TemplateID: "GERMAN_TEMPLATE",
					Path:       "/testdata/languages/de/vocabulary/s.md",
					Position:   "Spaziergang",
				},
				{
					Fields: map[string]string{
						"name": "Spiegel",
					},
					TemplateID: "GERMAN_TEMPLATE",
					Path:       "/testdata/languages/de/vocabulary/s.md",
					Position:   "Spiegel",
				},
			}},
		},
		{
			name: "should parse vocabulary and custom fields",
			path: "/testdata/full/example/a.md",
			config: config.VocabularyTemplate{
				TemplateID: "TEMPLATE_ID",
				ExamplesID: "EXAMPLES_ID",
				NotesID:    "NOTES_ID",
			},
			source: vocabularySource,
			want: Result{Cards: []Card{
				{
					Fields: map[string]string{
						"name":        "Spaziergang",
						"EXAMPLES_ID": "Wir haben nach dem Essen einen langen Spaziergang gemacht.",
						"NOTES_ID":    "Stroll, walk, promenade.",
					},
					TemplateID: "TEMPLATE_ID",
					Path:       "/testdata/full/example/a.md",
					Position:   "Spaziergang",
				},
				{
					Fields: map[string]string{
						"name": "Spiegel",
					},
					TemplateID: "TEMPLATE_ID",
					Path:       "/testdata/full/example/a.md",
					Position:   "Spiegel",
				},
				{
					Fields: map[string]string{
						"name": "First line can be a sentence.",
					},
					TemplateID: "TEMPLATE_ID",
					Path:       "/testdata/full/example/a.md",
					Position:   "Firstlinecanbeasentence",
				},
				{
					Fields: map[string]string{
						"name":        "Word",
						"EXAMPLES_ID": "Example without note for word.",
					},
					TemplateID: "TEMPLATE_ID",
					Path:       "/testdata/full/example/a.md",
					Position:   "Word",
				},
				{
					Fields: map[string]string{
						"name":     "AnotherWord",
						"NOTES_ID": "Notes without example.\n\nNotes can be multiline.",
					},
					TemplateID: "TEMPLATE_ID",
					Path:       "/testdata/full/example/a.md",
					Position:   "AnotherWord",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newVocabulary(tt.config).parse(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
