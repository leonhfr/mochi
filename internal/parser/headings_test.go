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

func Test_headings_convert(t *testing.T) {
	tests := []struct {
		name     string
		maxLevel int
		path     string
		source   string
		want     Result
	}{
		{
			name:     "empty file",
			maxLevel: 1,
			path:     "/Empty.md",
			want:     Result{Deck: "Empty"},
		},
		{
			name:     "simple level 1",
			maxLevel: 1,
			path:     "/Headings.md",
			source:   "# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n\n# Heading 2\n\nContent 3.\n",
			want: Result{Deck: "Headings", Cards: []Card{
				headingsCard{
					name:     "Headings > Heading 1",
					content:  "Headings > Heading 1\n\n# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n",
					path:     "/Headings.md",
					position: "Headingsmd0000",
				},
				headingsCard{
					name:     "Headings > Heading 2",
					content:  "Headings > Heading 2\n\n# Heading 2\n\nContent 3.\n",
					path:     "/Headings.md",
					position: "Headingsmd0001",
				},
			}},
		},
		{
			name:     "level 1 only headers",
			maxLevel: 1,
			path:     "/Headings.md",
			source:   "# Title 1\n\n# Title 2\n\n# Title 3\n\n# Title 4\n\n# Title 5\n",
			want:     Result{Deck: "Headings", Cards: []Card{}},
		},
		{
			name:     "level 1",
			maxLevel: 1,
			path:     "/Headings.md",
			source:   headingSource,
			want: Result{Deck: "Headings", Cards: []Card{
				headingsCard{
					name:     "Headings > Heading 1",
					content:  "Headings > Heading 1\n\n# Heading 1\n\nSome text here.\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n\n## Heading 1.2\n\nAnother content.\n",
					path:     "/Headings.md",
					position: "Headingsmd0000",
				},
				headingsCard{
					name:     "Headings > Heading 2",
					content:  "Headings > Heading 2\n\n# Heading 2\n\nCard card card.\n",
					path:     "/Headings.md",
					position: "Headingsmd0001",
				},
				headingsCard{
					name:     "Headings > Heading 3",
					content:  "Headings > Heading 3\n\n# Heading 3\n\n## Heading 3.1\n\nMore card content.\n",
					path:     "/Headings.md",
					position: "Headingsmd0002",
				},
			}},
		},
		{
			name:     "level 2",
			maxLevel: 2,
			path:     "/Headings.md",
			source:   headingSource,
			want: Result{Deck: "Headings", Cards: []Card{
				headingsCard{
					name:     "Headings > Heading 1",
					content:  "Headings > Heading 1\n\n# Heading 1\n\nSome text here.\n",
					path:     "/Headings.md",
					position: "Headingsmd0000",
				},
				headingsCard{
					name:     "Headings > Heading 1 > Heading 1.1",
					content:  "Headings > Heading 1 > Heading 1.1\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n",
					path:     "/Headings.md",
					position: "Headingsmd0001",
				},
				headingsCard{
					name:     "Headings > Heading 1 > Heading 1.2",
					content:  "Headings > Heading 1 > Heading 1.2\n\n## Heading 1.2\n\nAnother content.\n",
					path:     "/Headings.md",
					position: "Headingsmd0002",
				},
				headingsCard{
					name:     "Headings > Heading 2",
					content:  "Headings > Heading 2\n\n# Heading 2\n\nCard card card.\n",
					path:     "/Headings.md",
					position: "Headingsmd0003",
				},
				headingsCard{
					name:     "Headings > Heading 3 > Heading 3.1",
					content:  "Headings > Heading 3 > Heading 3.1\n\n## Heading 3.1\n\nMore card content.\n",
					path:     "/Headings.md",
					position: "Headingsmd0004",
				},
			}},
		},
		{
			name:     "level 2 with skip level",
			maxLevel: 2,
			path:     "/Headings.md",
			source:   "## Title 1\n\nContent 1.\n\n## Title 2\n\n### Title 2.1\n\nContent 1.\n",
			want: Result{Deck: "Headings", Cards: []Card{
				headingsCard{
					name:     "Headings > Title 1",
					content:  "Headings > Title 1\n\n## Title 1\n\nContent 1.\n",
					path:     "/Headings.md",
					position: "Headingsmd0000",
				},
				headingsCard{
					name:     "Headings > Title 2",
					content:  "Headings > Title 2\n\n## Title 2\n\n### Title 2.1\n\nContent 1.\n",
					path:     "/Headings.md",
					position: "Headingsmd0001",
				},
			}},
		},
		{
			name:     "images",
			maxLevel: 1,
			path:     "/subdirectory/Images.md",
			source:   "# Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n\n# Heading 2\n\n![Example 2](images/example-2.png)\n",
			want: Result{Deck: "Images", Cards: []Card{
				headingsCard{
					name:     "Images > Heading 1",
					content:  "Images > Heading 1\n\n# Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n",
					path:     "/subdirectory/Images.md",
					images:   []Image{{Destination: "../images/example-1.png", AltText: "Example 1"}},
					position: "Imagesmd0000",
				},
				headingsCard{
					name:     "Images > Heading 2",
					content:  "Images > Heading 2\n\n# Heading 2\n\n![Example 2](images/example-2.png)\n",
					path:     "/subdirectory/Images.md",
					images:   []Image{{Destination: "images/example-2.png", AltText: "Example 2"}},
					position: "Imagesmd0001",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHeadings(tt.maxLevel).convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
