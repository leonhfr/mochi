package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Headings_Convert(t *testing.T) {
	tests := []struct {
		name   string
		level  int
		path   string
		source string
		cards  []Card
		err    error
	}{
		{
			name:   "valid",
			level:  1,
			path:   "/headings.md",
			source: "<!-- Comment. -->\n\n# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n\n# Heading 2\n\nContent 3.\n",
			cards: []Card{
				{
					Name:     "Heading 1",
					Content:  "# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n",
					Filename: "headings.md",
				},
				{
					Name:     "Heading 2",
					Content:  "# Heading 2\n\nContent 3.\n",
					Filename: "headings.md",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHeadings(tt.level).convert(tt.path, []byte(tt.source))
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.cards, got)
		})
	}
}
