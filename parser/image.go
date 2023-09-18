package parser

import (
	//nolint:gosec
	"crypto/md5"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

const fileNameLength = 16

type Image struct {
	Destination string // original destination
	FileName    string // <md5 of path relative to root>
	Extension   string // ext
	ContentType string // mime type
	AltText     string
}

var mimeTypes = map[string]string{
	"bmp":  "image/bmp",
	"gif":  "image/gif",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpg",
	"png":  "image/png",
	"webp": "image/webp",
}

func newImage(path, destination, altText string) (string, Image) {
	if _, err := url.ParseRequestURI(destination); err == nil {
		return "", Image{}
	}

	if filepath.IsAbs(destination) {
		return "", Image{}
	}

	ext := strings.TrimLeft(filepath.Ext(destination), ".")
	mime, ok := mimeTypes[ext]
	if !ok {
		return "", Image{}
	}

	absPath := filepath.Join(filepath.Dir(path), destination)
	//nolint:gosec
	hash := fmt.Sprintf("%x", md5.Sum([]byte(absPath)))

	return absPath, Image{
		Destination: destination,
		FileName:    hash[:fileNameLength],
		Extension:   ext,
		AltText:     altText,
		ContentType: mime,
	}
}

func replaceImages(source string, images map[string]Image) string {
	for _, image := range images {
		from := fmt.Sprintf("![%s](%s)", image.AltText, image.Destination)
		to := fmt.Sprintf("![%s](@media/%s.%s)", image.AltText, image.FileName, image.Extension)
		source = strings.ReplaceAll(source, from, to)
	}
	return source
}
