package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Headings_Convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
		cards  []Card
		err    error
	}{
		{
			"valid",
			"/headings.md",
			"<!-- Comment. -->\n\n# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n\n# Heading 2\n\nContent 3.\n",
			[]Card{
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
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHeadings().convert(tt.path, []byte(tt.source))
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.cards, got)
		})
	}
}
