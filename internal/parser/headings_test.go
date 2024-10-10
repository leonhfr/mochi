package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var headingSource = `# Heading 1

Some text here.

## Heading 1.1

### Heading 1.1.1

Actual content.

More content.

## Heading 1.2

Another content.

# Heading 2

Card card card.

# Heading 3

## Heading 3.1

More card content.
`

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
			name:   "simple level 1",
			level:  1,
			path:   "/headings.md",
			source: "# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n\n# Heading 2\n\nContent 3.\n",
			cards: []Card{
				{
					Name:     "Heading 1", // TODO
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
		{
			name:   "level 1",
			level:  1,
			path:   "/headings.md",
			source: headingSource,
			cards: []Card{
				{
					Name:     "Heading 1", // todo
					Content:  "# Heading 1\n\nSome text here.\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n\n## Heading 1.2\n\nAnother content.\n",
					Filename: "headings.md",
				},
				{
					Name:     "Heading 2",
					Content:  "# Heading 2\n\nCard card card.\n",
					Filename: "headings.md",
				},
				{
					Name:     "Heading 3",
					Content:  "# Heading 3\n\n## Heading 3.1\n\nMore card content.\n",
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
