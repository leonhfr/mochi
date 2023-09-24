package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Note_Convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		source string
		want   []Card
	}{
		{
			"comment",
			"/note.md",
			"<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
			[]Card{
				{
					Name:    "Noun",
					Content: "<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
					Fields:  map[string]string{},
					Images:  map[string]Image{},
				},
			},
		},
		{
			"front matter",
			"/note.md",
			"---\nfoo: bar\n---\n\n<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
			[]Card{
				{
					Name:    "Noun",
					Content: "<!-- Comment. -->\n\n# Noun\n\n## Gender\n\nSome stuff about genders.\n\n- der\n- die\n- das\n",
					Fields:  map[string]string{},
					Images:  map[string]Image{},
				},
			},
		},
		{
			"images",
			"/dir/note.md",
			"# Noun\n\n## Gender\n\n![Example 1](../images/example-1.png)\n\nInline image: ![Example 2](./example-2.png)",
			[]Card{
				{
					Name:    "Noun",
					Content: "# Noun\n\n## Gender\n\n![Example 1](@media/db7ab4bbe96b326a.png)\n\nInline image: ![Example 2](@media/bd2c42f53f241cba.png)\n",
					Fields:  map[string]string{},
					Images: map[string]Image{
						"/images/example-1.png": {
							Destination: "../images/example-1.png",
							FileName:    "db7ab4bbe96b326a",
							Extension:   "png",
							ContentType: "image/png",
							AltText:     "Example 1",
						},
						"/dir/example-2.png": {
							Destination: "./example-2.png",
							FileName:    "bd2c42f53f241cba",
							Extension:   "png",
							ContentType: "image/png",
							AltText:     "Example 2",
						},
					},
				},
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNote().Convert(tt.path, []byte(tt.source))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
