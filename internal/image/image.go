package image

import (
	//nolint:gosec
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"
)

const fileNameLength = 16

var mimeTypes = map[string]string{
	"bmp":  "image/bmp",
	"gif":  "image/gif",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpg",
	"png":  "image/png",
	"webp": "image/webp",
}

// FileCheck is the interface implemented to check file existence.
type FileCheck interface {
	Exists(path string) bool
}

// Image contains the data of one image.
type Image struct {
	Filename    string // md5 of path relative to root
	Extension   string
	MimeType    string
	Destination string // original destination
	AltText     string
}

// Parsed contains the data of a parsed image.
type Parsed struct {
	Destination string
	AltText     string
}

// New creates a new Image.
func New(fc FileCheck, path string, parsed Parsed) (string, Image, bool) {
	dirPath := filepath.Join("./", filepath.Dir(path))
	absPath := filepath.Join(dirPath, parsed.Destination)
	if !fc.Exists(absPath) {
		return "", Image{}, false
	}

	if filepath.IsAbs(parsed.Destination) {
		absPath = parsed.Destination
	}

	ext := strings.TrimLeft(filepath.Ext(parsed.Destination), ".")
	mime, ok := mimeTypes[ext]
	if !ok {
		return "", Image{}, false
	}

	//nolint:gosec
	hash := fmt.Sprintf("%x", md5.Sum([]byte(absPath)))
	return absPath, Image{
		Filename:    hash[:fileNameLength],
		Extension:   ext,
		MimeType:    mime,
		Destination: parsed.Destination,
		AltText:     parsed.AltText,
	}, true
}

// NewMap create a new map from parsed images.
func NewMap(fc FileCheck, path string, parsed []Parsed) map[string]Image {
	images := make(map[string]Image)
	for _, p := range parsed {
		if absPath, image, ok := New(fc, path, p); ok {
			images[absPath] = image
		}
	}
	return images
}

// Replace replaces images link in the Markdown source to mochi Markdown.
func Replace(images map[string]Image, source string) string {
	for _, image := range images {
		from := fmt.Sprintf("![%s](%s)", image.AltText, image.Destination)
		to := fmt.Sprintf("![%s](@media/%s.%s)", image.AltText, image.Filename, image.Extension)
		source = strings.ReplaceAll(source, from, to)
	}
	return source
}
