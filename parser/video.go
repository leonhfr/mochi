package parser

import (
	"regexp"
	"strings"
)

const (
	iframePrefix  = `<iframe src="`
	iframeOptions = `?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0" frameborder="0" loading="lazy" gesture="media" allow="autoplay; fullscreen" allowautoplay="true" allowfullscreen="true" style="`
	iframeStyle   = `aspect-ratio:16/9;height:100%;width:100%;`
	iframeSuffix  = `"></iframe>`
)

func replaceVideos(content string) string {
	video := regexp.MustCompile(`(?m)^\[@video\]\((.+)\)$`)
	return video.ReplaceAllString(content, formatVideo("$1"))
}

func formatVideo(dest string) string {
	return strings.Join([]string{
		iframePrefix,
		dest,
		iframeOptions,
		iframeStyle,
		iframeSuffix,
	}, "")
}
