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

// Map contains the data of all images in the file.
type Map struct {
	fileCheck FileCheck
	dirPath   string
	images    map[string]Image
}

// Image contains the data of one image.
type Image struct {
	Filename    string // md5 of path relative to root
	Extension   string
	MimeType    string
	destination string // original destination
	altText     string
}

// New returns a new Images map.
func New(fileCheck FileCheck, path string) Map {
	return Map{
		fileCheck: fileCheck,
		dirPath:   fmt.Sprintf(".%s", filepath.Dir(path)),
		images:    make(map[string]Image),
	}
}

// Add adds an image.
func (i *Map) Add(destination, altText string) {
	absPath := filepath.Join(i.dirPath, destination)
	if !i.fileCheck.Exists(absPath) {
		return
	}

	if filepath.IsAbs(destination) {
		absPath = destination
	}

	if _, ok := i.images[absPath]; ok {
		return
	}

	ext := strings.TrimLeft(filepath.Ext(destination), ".")
	mime, ok := mimeTypes[ext]
	if !ok {
		return
	}

	//nolint:gosec
	hash := fmt.Sprintf("%x", md5.Sum([]byte(absPath)))
	i.images[absPath] = Image{
		Filename:    hash[:fileNameLength],
		Extension:   ext,
		MimeType:    mime,
		destination: destination,
		altText:     altText,
	}
}

// Replace replaces images link in the Markdown source to mochi Markdown.
func (i *Map) Replace(source string) string {
	for _, image := range i.images {
		from := fmt.Sprintf("![%s](%s)", image.altText, image.destination)
		to := fmt.Sprintf("![%s](@media/%s.%s)", image.altText, image.Filename, image.Extension)
		source = strings.ReplaceAll(source, from, to)
	}
	return source
}

// Images returns the images map.
func (i *Map) Images() map[string]Image {
	return i.images
}
