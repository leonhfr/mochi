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

func Test_headings_parse(t *testing.T) {
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
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 1</details>\n\n# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n",
					Fields:   nameFields("Headings > Heading 1"),
					Path:     "/Headings.md",
					Position: "Headingsmd0000",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 2</details>\n\n# Heading 2\n\nContent 3.\n",
					Fields:   nameFields("Headings > Heading 2"),
					Path:     "/Headings.md",
					Position: "Headingsmd0001",
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
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 1</details>\n\n# Heading 1\n\nSome text here.\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n\n## Heading 1.2\n\nAnother content.\n",
					Fields:   nameFields("Headings > Heading 1"),
					Path:     "/Headings.md",
					Position: "Headingsmd0000",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 2</details>\n\n# Heading 2\n\nCard card card.\n",
					Fields:   nameFields("Headings > Heading 2"),
					Path:     "/Headings.md",
					Position: "Headingsmd0001",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 3</details>\n\n# Heading 3\n\n## Heading 3.1\n\nMore card content.\n",
					Fields:   nameFields("Headings > Heading 3"),
					Path:     "/Headings.md",
					Position: "Headingsmd0002",
				},
			}},
		},
		{
			name:     "level 2",
			maxLevel: 2,
			path:     "/Headings.md",
			source:   headingSource,
			want: Result{Deck: "Headings", Cards: []Card{
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 1</details>\n\n# Heading 1\n\nSome text here.\n",
					Fields:   nameFields("Headings > Heading 1"),
					Path:     "/Headings.md",
					Position: "Headingsmd0000",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 1 > Heading 1.1</details>\n\n## Heading 1.1\n\n### Heading 1.1.1\n\nActual content.\n\nMore content.\n",
					Fields:   nameFields("Headings > Heading 1 > Heading 1.1"),
					Path:     "/Headings.md",
					Position: "Headingsmd0001",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 1 > Heading 1.2</details>\n\n## Heading 1.2\n\nAnother content.\n",
					Fields:   nameFields("Headings > Heading 1 > Heading 1.2"),
					Path:     "/Headings.md",
					Position: "Headingsmd0002",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 2</details>\n\n# Heading 2\n\nCard card card.\n",
					Fields:   nameFields("Headings > Heading 2"),
					Path:     "/Headings.md",
					Position: "Headingsmd0003",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Heading 3 > Heading 3.1</details>\n\n## Heading 3.1\n\nMore card content.\n",
					Fields:   nameFields("Headings > Heading 3 > Heading 3.1"),
					Path:     "/Headings.md",
					Position: "Headingsmd0004",
				},
			}},
		},
		{
			name:     "level 2 with skip level",
			maxLevel: 2,
			path:     "/Headings.md",
			source:   "## Title 1\n\nContent 1.\n\n## Title 2\n\n### Title 2.1\n\nContent 1.\n",
			want: Result{Deck: "Headings", Cards: []Card{
				{
					Content:  "<details><summary>Headings</summary>Headings > Title 1</details>\n\n## Title 1\n\nContent 1.\n",
					Fields:   nameFields("Headings > Title 1"),
					Path:     "/Headings.md",
					Position: "Headingsmd0000",
				},
				{
					Content:  "<details><summary>Headings</summary>Headings > Title 2</details>\n\n## Title 2\n\n### Title 2.1\n\nContent 1.\n",
					Fields:   nameFields("Headings > Title 2"),
					Path:     "/Headings.md",
					Position: "Headingsmd0001",
				},
			}},
		},
		{
			name:     "images",
			maxLevel: 1,
			path:     "/subdirectory/Images.md",
			source:   "# Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n\n# Heading 2\n\n![Example 2](images/example-2.png)\n",
			want: Result{Deck: "Images", Cards: []Card{
				{
					Content:  "<details><summary>Headings</summary>Images > Heading 1</details>\n\n# Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n",
					Fields:   nameFields("Images > Heading 1"),
					Path:     "/subdirectory/Images.md",
					Position: "Imagesmd0000",
				},
				{
					Content:  "<details><summary>Headings</summary>Images > Heading 2</details>\n\n# Heading 2\n\n![Example 2](images/example-2.png)\n",
					Fields:   nameFields("Images > Heading 2"),
					Path:     "/subdirectory/Images.md",
					Position: "Imagesmd0001",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHeadings(tt.maxLevel).parse(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
