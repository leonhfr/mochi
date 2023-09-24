package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replaceVideos(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			"videos",
			"[@video](https://www.youtube-nocookie.com/embed/FIRST_ID)\n\n[@video](https://www.youtube-nocookie.com/embed/SECOND_ID)",
			"<iframe src=\"https://www.youtube-nocookie.com/embed/FIRST_ID?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0\" frameborder=\"0\" loading=\"lazy\" gesture=\"media\" allow=\"autoplay; fullscreen\" allowautoplay=\"true\" allowfullscreen=\"true\" style=\"aspect-ratio:16/9;height:100%;width:100%;\"></iframe>\n\n<iframe src=\"https://www.youtube-nocookie.com/embed/SECOND_ID?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0\" frameborder=\"0\" loading=\"lazy\" gesture=\"media\" allow=\"autoplay; fullscreen\" allowautoplay=\"true\" allowfullscreen=\"true\" style=\"aspect-ratio:16/9;height:100%;width:100%;\"></iframe>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceVideos(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}
