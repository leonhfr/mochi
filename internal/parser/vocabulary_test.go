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

func Test_vocabulary_convert(t *testing.T) {
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
				vocabularyCard{
					config: config.VocabularyTemplate{TemplateID: "GERMAN_TEMPLATE"},
					word:   "Spaziergang",
					path:   "/testdata/languages/de/vocabulary/s.md",
				},
				vocabularyCard{
					config: config.VocabularyTemplate{TemplateID: "GERMAN_TEMPLATE"},
					word:   "Spiegel",
					path:   "/testdata/languages/de/vocabulary/s.md",
				},
			}},
		},
		{
			name:   "should parse vocabulary and custom fields",
			path:   "/testdata/full/example/a.md",
			config: config.VocabularyTemplate{TemplateID: "FULL_EXAMPLE"},
			source: vocabularySource,
			want: Result{Cards: []Card{
				vocabularyCard{
					config:   config.VocabularyTemplate{TemplateID: "FULL_EXAMPLE"},
					word:     "Spaziergang",
					examples: []string{"Wir haben nach dem Essen einen langen Spaziergang gemacht."},
					notes:    []string{"Stroll, walk, promenade."},
					path:     "/testdata/full/example/a.md",
				},
				vocabularyCard{
					config: config.VocabularyTemplate{TemplateID: "FULL_EXAMPLE"},
					word:   "Spiegel",
					path:   "/testdata/full/example/a.md",
				},
				vocabularyCard{
					config: config.VocabularyTemplate{TemplateID: "FULL_EXAMPLE"},
					word:   "First line can be a sentence.",
					path:   "/testdata/full/example/a.md",
				},
				vocabularyCard{
					config:   config.VocabularyTemplate{TemplateID: "FULL_EXAMPLE"},
					word:     "Word",
					examples: []string{"Example without note for word."},
					path:     "/testdata/full/example/a.md",
				},
				vocabularyCard{
					config: config.VocabularyTemplate{TemplateID: "FULL_EXAMPLE"},
					word:   "AnotherWord",
					notes:  []string{"Notes without example.", "Notes can be multiline."},
					path:   "/testdata/full/example/a.md",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newVocabulary(tt.config).convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
