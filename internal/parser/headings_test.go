package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/image"
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
		name       string
		maxLevel   int
		path       string
		source     string
		fileChecks map[string]bool
		want       []Card
	}{
		{
			name:     "simple level 1",
			maxLevel: 1,
			path:     "/Headings.md",
			source:   "# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n\n# Heading 2\n\nContent 3.\n",
			want: []Card{
				{
					Name:     "Headings | Heading 1",
					Content:  "Headings | Heading 1\n\n# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
				},
				{
					Name:     "Headings | Heading 2",
					Content:  "Headings | Heading 2\n\n# Heading 2\n\nContent 3.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    1,
				},
			},
		},
		{
			name:     "level 1",
			maxLevel: 1,
			path:     "/Headings.md",
			source:   headingSource,
			want: []Card{
				{
					Name:     "Headings | Heading 1",
					Content:  "Headings | Heading 1\n\n# Heading 1\n\nSome text here.\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n\n## Heading 1.2\n\nAnother content.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
				},
				{
					Name:     "Headings | Heading 2",
					Content:  "Headings | Heading 2\n\n# Heading 2\n\nCard card card.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    1,
				},
				{
					Name:     "Headings | Heading 3",
					Content:  "Headings | Heading 3\n\n# Heading 3\n\n## Heading 3.1\n\nMore card content.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    2,
				},
			},
		},
		{
			name:     "level 2",
			maxLevel: 2,
			path:     "/Headings.md",
			source:   headingSource,
			want: []Card{
				{
					Name:     "Headings | Heading 1",
					Content:  "Headings | Heading 1\n\n# Heading 1\n\nSome text here.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
				},
				{
					Name:     "Headings | Heading 1 | Heading 1.1",
					Content:  "Headings | Heading 1 | Heading 1.1\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    1,
				},
				{
					Name:     "Headings | Heading 1 | Heading 1.2",
					Content:  "Headings | Heading 1 | Heading 1.2\n\n## Heading 1.2\n\nAnother content.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    2,
				},
				{
					Name:     "Headings | Heading 2",
					Content:  "Headings | Heading 2\n\n# Heading 2\n\nCard card card.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    3,
				},
				{
					Name:     "Headings | Heading 3 | Heading 3.1",
					Content:  "Headings | Heading 3 | Heading 3.1\n\n# Heading 3\n\n## Heading 3.1\n\nMore card content.\n",
					Filename: "Headings.md",
					Images:   map[string]image.Image{},
					Index:    4,
				},
			},
		},
		{
			name:     "images",
			maxLevel: 1,
			path:     "/subdirectory/Images.md",
			source:   "# Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n\n# Heading 2\n\n![Example 2](images/example-2.png)\n",
			fileChecks: map[string]bool{
				"/images/example-1.png":              true,
				"/subdirectory/images/example-2.png": true,
			},
			want: []Card{
				{
					Name:     "Images | Heading 1",
					Content:  "Images | Heading 1\n\n# Heading 1\n\nContent 1.\n\n![Example 1](@media/db7ab4bbe96b326a.png)\n",
					Filename: "Images.md",
					Images:   map[string]image.Image{"/images/example-1.png": {Filename: "db7ab4bbe96b326a", Extension: "png", MimeType: "image/png", Destination: "../images/example-1.png", AltText: "Example 1"}},
				},
				{
					Name:     "Images | Heading 2",
					Content:  "Images | Heading 2\n\n# Heading 2\n\n![Example 2](@media/01a5479a1f430b25.png)\n",
					Filename: "Images.md",
					Images:   map[string]image.Image{"/subdirectory/images/example-2.png": {Filename: "01a5479a1f430b25", Extension: "png", MimeType: "image/png", Destination: "images/example-2.png", AltText: "Example 2"}},
					Index:    1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newMockFileChecker(tt.fileChecks)
			got, err := newHeadings(fc, tt.maxLevel).convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			fc.AssertExpectations(t)
		})
	}
}
