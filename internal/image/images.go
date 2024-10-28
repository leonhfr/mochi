package image

import (
	"fmt"
	"strings"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

// Images contains the images data.
type Images []Image

// New creates a new image list from parsed images.
func New(reader Reader, path string, parsed []parser.Image) Images {
	images := []Image{}
	for _, p := range parsed {
		if image, ok := newImage(reader, path, p); ok {
			images = append(images, image)
		}
	}
	return images
}

// Attachments returns the list of attachments.
func (images Images) Attachments() []mochi.DeprecatedAttachment {
	attachments := make([]mochi.DeprecatedAttachment, len(images))
	for i, image := range images {
		attachments[i] = image.attachment
	}
	return attachments
}

// HashMap returns the map of [abs path]hash.
func (images Images) HashMap() map[string]string {
	hashMap := make(map[string]string)
	for _, image := range images {
		hashMap[image.Path] = image.Hash
	}
	return hashMap
}

// Replace replaces images link in the Markdown source to mochi Markdown.
func (images Images) Replace(content string) string {
	for _, image := range images {
		from := fmt.Sprintf("![%s](%s)", image.altText, image.destination)
		to := fmt.Sprintf("![%s](@media/%s)", image.altText, image.attachment.FileName)
		content = strings.ReplaceAll(content, from, to)
	}
	return content
}
