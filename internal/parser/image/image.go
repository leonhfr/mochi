package image

import (
	//nolint:gosec
	"crypto/md5"
	"fmt"
	"net/url"
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

// Map contains the data of all images in the file.
type Map struct {
	dirPath string
	// absolutePath string // path to image from deck root
	images map[string]Image
}

// Image contains the data of one image.
type Image struct {
	fileName    string // md5 of path relative to root
	destination string // original destination
	extension   string
	mimeType    string
	altText     string
}

// New returns a new Images map.
func New(path string) Map {
	return Map{
		dirPath: filepath.Dir(path),
		images:  make(map[string]Image),
	}
}

// Add adds an image.
func (i *Map) Add(destination, altText string) {
	absPath := filepath.Join(i.dirPath, destination)
	if _, ok := i.images[absPath]; ok {
		return
	}

	if _, err := url.ParseRequestURI(destination); err == nil {
		return
	}

	if filepath.IsAbs(destination) {
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
		fileName:    hash[:fileNameLength],
		destination: destination,
		extension:   ext,
		mimeType:    mime,
		altText:     altText,
	}
}

// Replace replaces images link in the Markdown source to mochi Markdown.
func (i *Map) Replace(source string) string {
	for _, image := range i.images {
		from := fmt.Sprintf("![%s](%s)", image.altText, image.destination)
		to := fmt.Sprintf("![%s](@media/%s.%s)", image.altText, image.fileName, image.extension)
		source = strings.ReplaceAll(source, from, to)
	}
	return source
}
