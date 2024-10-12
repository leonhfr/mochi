package image

import (
	"bytes"
	//nolint:gosec
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/leonhfr/mochi/mochi"
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
	filename    string // md5 of path relative to root
	destination string // original destination
	extension   string
	mimeType    string
	altText     string
}

// New returns a new Images map.
func New(path string) Map {
	return Map{
		dirPath: fmt.Sprintf(".%s", filepath.Dir(path)),
		images:  make(map[string]Image),
	}
}

// Add adds an image.
func (i *Map) Add(destination, altText string) {
	absPath := filepath.Join(i.dirPath, destination)
	if _, err := os.Stat(absPath); err != nil {
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
		filename:    hash[:fileNameLength],
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
		to := fmt.Sprintf("![%s](@media/%s.%s)", image.altText, image.filename, image.extension)
		source = strings.ReplaceAll(source, from, to)
	}
	return source
}

// Reader represents the interface to read image files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Attachment contains the data of a mochi attachment.
type Attachment struct {
	Mochi mochi.Attachment
	Hash  string
	Path  string
}

// Attachments converts the images map to mochi attachments.
func (i *Map) Attachments(r Reader) ([]Attachment, error) {
	attachments := make([]Attachment, 0, len(i.images))
	for path, image := range i.images {
		hash, attachment, err := newMochiAttachment(r, path, image)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, Attachment{
			Mochi: attachment,
			Hash:  hash,
			Path:  path,
		})
	}
	return attachments, nil
}

func newMochiAttachment(r Reader, path string, image Image) (string, mochi.Attachment, error) {
	file, err := r.Read(path)
	if err != nil {
		return "", mochi.Attachment{}, err
	}
	defer file.Close()

	bytes := bytes.NewBuffer(nil)
	base64Encoder := base64.NewEncoder(base64.StdEncoding, bytes)
	defer base64Encoder.Close()

	//nolint:gosec
	hashEncoder := md5.New()
	tee := io.TeeReader(file, hashEncoder)
	if _, err := io.Copy(base64Encoder, tee); err != nil {
		return "", mochi.Attachment{}, err
	}

	hash := fmt.Sprintf("%x", hashEncoder.Sum(nil))
	return hash, mochi.Attachment{
		FileName:    fmt.Sprintf("%s.%s", image.filename, image.extension),
		ContentType: image.mimeType,
		Data:        bytes.String(),
	}, nil
}
