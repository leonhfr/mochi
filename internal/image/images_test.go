package image

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

func Test_New(t *testing.T) {
	path := "/testdata/Markdown.md"
	tests := []struct {
		name   string
		calls  []testRead
		parsed []parser.Image
		want   Images
	}{
		{
			name: "should return the expected slice",
			calls: []testRead{
				{
					path:    "/testdata/scream.png",
					content: "IMAGE CONTENT",
				},
				{
					path: "/testdata/unknown.png",
					err:  fs.ErrNotExist,
				},
			},
			parsed: []parser.Image{
				{
					Destination: "unknown.png",
					AltText:     "alt text",
				},
				{
					Destination: "scream.png",
					AltText:     "alt text",
				},
			},
			want: Images{
				{
					attachment: mochi.Attachment{
						FileName:    "22abb8f07c02970e.png",
						ContentType: "image/png",
						Data:        "SU1BR0UgQ09OVEVO",
					},
					Hash:        "1923784bcb1663bbbd9efd9765c36382",
					Path:        "/testdata/scream.png",
					destination: "scream.png",
					altText:     "alt text",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader(tt.calls)
			got := New(r, path, tt.parsed)
			assert.Equal(t, tt.want, got)
			r.AssertExpectations(t)
		})
	}
}

func Test_Images_Replace(t *testing.T) {
	tests := []struct {
		name    string
		images  Images
		content string
		want    string
	}{
		{
			name: "should replace the images",
			images: Images{
				{
					attachment:  mochi.Attachment{FileName: "scream_hash.png"},
					destination: "./scream.png",
					altText:     "Scream",
				},
				{
					attachment:  mochi.Attachment{FileName: "constellations_hash.jpg"},
					destination: "./constellations.jpg",
				},
			},
			content: "![Scream](./scream.png)\n![](./constellations.jpg)",
			want:    "![Scream](@media/scream_hash.png)\n![](@media/constellations_hash.jpg)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.images.Replace(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}
