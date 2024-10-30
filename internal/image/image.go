package image

import (
	"bytes"
	//nolint:gosec
	"crypto/md5"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/leonhfr/mochi/internal/parser"
)

const fileNameLength = 16

// Reader represents the interface to read files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Image contains the data for an Image.
type Image struct {
	Bytes       []byte // image bytes
	Filename    string // filename: [md5 of path].ext
	Hash        string // md5 of contents
	Path        string // absolute path to file
	destination string // original destination
	altText     string // original alt text
}

// newImage creates a new image.
func newImage(reader Reader, path string, parsed parser.Image) (Image, bool) {
	absPath := filepath.Join(filepath.Dir(path), parsed.Destination)
	hash, bytes, err := readImage(reader, absPath)
	if err != nil {
		return Image{}, false
	}

	extension := getExtension(parsed.Destination)
	pathHash := getPathHash(absPath)
	filename := getFilename(pathHash, extension)

	return Image{
		Bytes:       bytes,
		Filename:    filename,
		Path:        absPath,
		Hash:        hash,
		destination: parsed.Destination,
		altText:     parsed.AltText,
	}, true
}

func readImage(reader Reader, absPath string) (string, []byte, error) {
	file, err := reader.Read(absPath)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	bytes := bytes.NewBuffer(nil)

	//nolint:gosec
	hashEncoder := md5.New()
	tee := io.TeeReader(file, hashEncoder)
	if _, err := io.Copy(bytes, tee); err != nil {
		return "", nil, err
	}

	hash := fmt.Sprintf("%x", hashEncoder.Sum(nil))
	return hash, bytes.Bytes(), nil
}

func getExtension(destination string) string {
	return strings.TrimLeft(filepath.Ext(destination), ".")
}

func getPathHash(absPath string) string {
	//nolint:gosec
	return fmt.Sprintf("%x", md5.Sum([]byte(absPath)))
}

func getFilename(pathHash, extension string) string {
	shortHash := pathHash[:fileNameLength]
	return fmt.Sprintf("%s.%s", shortHash, extension)
}
