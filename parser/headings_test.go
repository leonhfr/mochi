package parser

import (
	"errors"
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
					Name:    "Heading 1",
					Content: "# Heading 1\n\nContent 1.\n\n## Subtitle\n\nContent 2.\n",
					Fields:  map[string]string{},
					Images:  map[string]Image{},
				},
				{
					Name:    "Heading 2",
					Content: "# Heading 2\n\nContent 3.\n",
					Fields:  map[string]string{},
					Images:  map[string]Image{},
				},
			},
			nil,
		},
		{
			"images",
			"/subdirectory/headings.md",
			"---\nfoo: bar\n---\n\n# Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n\n# Heading 2\n\n![Example 2](images/example-2.png)\n",
			[]Card{
				{
					Name:    "Heading 1",
					Content: "# Heading 1\n\nContent 1.\n\n![Example 1](@media/db7ab4bbe96b326a.png)\n",
					Fields:  map[string]string{},
					Images: map[string]Image{
						"/images/example-1.png": {
							Destination: "../images/example-1.png",
							FileName:    "db7ab4bbe96b326a",
							Extension:   "png",
							ContentType: "image/png",
							AltText:     "Example 1",
						},
					},
				},
				{
					Name:    "Heading 2",
					Content: "# Heading 2\n\n![Example 2](@media/01a5479a1f430b25.png)\n",
					Fields:  map[string]string{},
					Images: map[string]Image{
						"/subdirectory/images/example-2.png": {
							Destination: "images/example-2.png",
							FileName:    "01a5479a1f430b25",
							Extension:   "png",
							ContentType: "image/png",
							AltText:     "Example 2",
						},
					},
				},
			},
			nil,
		},
		{
			"video",
			"/video.md",
			"# Video\n\n[@video](https://www.youtube-nocookie.com/embed/VIDEO_ID)\n",
			[]Card{
				{
					Name:    "Video",
					Content: "# Video\n\n<iframe src=\"https://www.youtube-nocookie.com/embed/VIDEO_ID?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0\" frameborder=\"0\" loading=\"lazy\" gesture=\"media\" allow=\"autoplay; fullscreen\" allowautoplay=\"true\" allowfullscreen=\"true\" style=\"aspect-ratio:16/9;height:100%;width:100%;\"></iframe>\n",
					Fields:  map[string]string{},
					Images:  map[string]Image{},
				},
			},
			nil,
		},
		{
			"malformed header",
			"/subdirectory/headings.md",
			"---\nfoo: bar\n---\n\n## Heading 1\n\nContent 1.\n\n![Example 1](../images/example-1.png)\n\n# Heading 2\n\nContent 2.\n",
			nil,
			errors.New("malformed card: /subdirectory/headings.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHeadings().Convert(tt.path, []byte(tt.source))
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.cards, got)
		})
	}
}
