package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ cardParser = &note{}

func Test_note_parse(t *testing.T) {
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
			want: Result{Cards: []Card{{
				Content: "# Lorem ipsum\n\nParagraph.\n",
				Fields:  nameFields("Lorem ipsum"),
				Images:  []Image{},
				Path:    "/testdata/lorem-ipsum/Lorem ipsum.md",
			}}},
		},
		{
			name:   "images",
			path:   "/testdata/lorem-ipsum/Lorem ipsum.md",
			source: "# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
			want: Result{Cards: []Card{{
				Content: "# Lorem ipsum\n\n# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
				Fields:  nameFields("Lorem ipsum"),
				Images: []Image{
					{Destination: "../images/example-1.png", AltText: "Example 1"},
					{Destination: "./example-2.png", AltText: "Example 2"},
				},
				Path: "/testdata/lorem-ipsum/Lorem ipsum.md",
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNote().parse(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
