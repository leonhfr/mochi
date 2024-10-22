package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/image"
)

var _ cardParser = &note{}

func Test_Note_Convert(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		source     string
		fileChecks map[string]bool
		want       []Card
	}{
		{
			name:   "simple note",
			path:   "/testdata/lorem-ipsum/Lorem ipsum.md",
			source: "# Title 1\nParagraph.\n",
			want: []Card{{
				Name:     "Lorem ipsum",
				Content:  "# Title 1\nParagraph.\n",
				Filename: "Lorem ipsum.md",
				Images:   map[string]image.Image{},
			}},
		},
		{
			name:   "images",
			path:   "/testdata/lorem-ipsum/Lorem ipsum.md",
			source: "# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
			fileChecks: map[string]bool{
				"/testdata/images/example-1.png":      true,
				"/testdata/lorem-ipsum/example-2.png": true,
			},
			want: []Card{{
				Name:     "Lorem ipsum",
				Content:  "# Noun\n\n## Gender\n\n![Example 1](@media/05142e92797cb245.png)\n\nInline image: ![Example 2](@media/2c36bee942023edc.png)",
				Filename: "Lorem ipsum.md",
				Images: map[string]image.Image{
					"/testdata/images/example-1.png":      {Filename: "05142e92797cb245", Extension: "png", MimeType: "image/png", Destination: "../images/example-1.png", AltText: "Example 1"},
					"/testdata/lorem-ipsum/example-2.png": {Filename: "2c36bee942023edc", Extension: "png", MimeType: "image/png", Destination: "./example-2.png", AltText: "Example 2"},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newMockFileChecker(tt.fileChecks)
			got, err := newNote().convert(fc, tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			fc.AssertExpectations(t)
		})
	}
}
