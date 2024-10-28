package image

import (
	"bytes"
	//nolint:gosec
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/leonhfr/mochi/internal/parser"
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

// Reader represents the interface to read files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Image contains the data for an Image.
type Image struct {
	attachment  mochi.DeprecatedAttachment // filename: [md5 of path].ext
	Hash        string                     // md5 of contents
	Path        string                     // absolute path to file
	destination string                     // original destination
	altText     string                     // original alt text
}

// newImage creates a new image.
func newImage(reader Reader, path string, parsed parser.Image) (Image, bool) {
	absPath := filepath.Join(filepath.Dir(path), parsed.Destination)
	hash, content, err := readImage(reader, absPath)
	if err != nil {
		return Image{}, false
	}

	extension := getExtension(parsed.Destination)
	mimeType, ok := getMimeType(extension)
	if !ok {
		return Image{}, false
	}

	pathHash := getPathHash(absPath)
	filename := getFilename(pathHash, extension)

	return Image{
		attachment: mochi.DeprecatedAttachment{
			FileName:    filename,
			ContentType: mimeType,
			Data:        content,
		},
		Path:        absPath,
		Hash:        hash,
		destination: parsed.Destination,
		altText:     parsed.AltText,
	}, true
}

func readImage(reader Reader, absPath string) (string, string, error) {
	file, err := reader.Read(absPath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	bytes := bytes.NewBuffer(nil)
	base64Encoder := base64.NewEncoder(base64.StdEncoding, bytes)
	defer base64Encoder.Close()

	//nolint:gosec
	hashEncoder := md5.New()
	tee := io.TeeReader(file, hashEncoder)
	if _, err := io.Copy(base64Encoder, tee); err != nil {
		return "", "", err
	}

	hash := fmt.Sprintf("%x", hashEncoder.Sum(nil))
	return hash, bytes.String(), nil
}

func getExtension(destination string) string {
	return strings.TrimLeft(filepath.Ext(destination), ".")
}

func getMimeType(extension string) (string, bool) {
	mime, ok := mimeTypes[extension]
	return mime, ok
}

func getPathHash(absPath string) string {
	//nolint:gosec
	return fmt.Sprintf("%x", md5.Sum([]byte(absPath)))
}

func getFilename(pathHash, extension string) string {
	shortHash := pathHash[:fileNameLength]
	return fmt.Sprintf("%s.%s", shortHash, extension)
}
