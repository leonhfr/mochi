package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ cardParser = &note{}

func Test_note_convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
		want   Result
	}{
		{
			name:   "simple note",
			path:   "/testdata/lorem-ipsum/Lorem ipsum.md",
			source: "Paragraph.\n",
			want: Result{Cards: []Card{noteCard{
				name:    "Lorem ipsum",
				content: "# Lorem ipsum\n\nParagraph.\n",
				path:    "/testdata/lorem-ipsum/Lorem ipsum.md",
				images:  []Image{},
			}}},
		},
		{
			name:   "images",
			path:   "/testdata/lorem-ipsum/Lorem ipsum.md",
			source: "# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
			want: Result{Cards: []Card{noteCard{
				name:    "Lorem ipsum",
				content: "# Lorem ipsum\n\n# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
				path:    "/testdata/lorem-ipsum/Lorem ipsum.md",
				images: []Image{
					{Destination: "../images/example-1.png", AltText: "Example 1"},
					{Destination: "./example-2.png", AltText: "Example 2"},
				},
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNote().convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
