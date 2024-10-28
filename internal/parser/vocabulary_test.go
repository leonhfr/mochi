package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_vocabulary_convert(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		templateID string
		source     string
		want       Result
	}{
		{
			name:       "should parse vocabulary",
			path:       "/testdata/languages/de/vocabulary/s.md",
			templateID: "GERMAN_TEMPLATE",
			source:     "<!-- Generated. -->\n\nSpaziergang\n\nSpiegel\n",
			want: Result{Cards: []Card{
				vocabularyCard{
					templateID: "GERMAN_TEMPLATE",
					word:       "Spaziergang",
					path:       "/testdata/languages/de/vocabulary/s.md",
				},
				vocabularyCard{
					templateID: "GERMAN_TEMPLATE",
					word:       "Spiegel",
					path:       "/testdata/languages/de/vocabulary/s.md",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newVocabulary(tt.templateID).convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
