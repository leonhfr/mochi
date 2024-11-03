package converter

import (
	"bytes"
	//nolint:gosec
	"crypto/md5"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

const fileNameLength = 16

// Attachment represents an attachment.
type Attachment struct {
	Bytes    []byte
	Filename string
}

func newAttachment(reader Reader, path, destination string) (Attachment, error) {
	absPath := filepath.Join(filepath.Dir(path), destination)
	bytes, err := readAttachment(reader, absPath)
	if err != nil {
		return Attachment{}, err
	}

	extension := getExtension(destination)
	pathHash := getPathHash(absPath)
	filename := getFilename(pathHash, extension)

	return Attachment{
		Bytes:    bytes,
		Filename: filename,
	}, nil
}

func (a Attachment) destination() []byte {
	return []byte(fmt.Sprintf("@media/%s", a.Filename))
}

func readAttachment(reader Reader, absPath string) ([]byte, error) {
	file, err := reader.Read(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes := bytes.NewBuffer(nil)

	if _, err := io.Copy(bytes, file); err != nil {
		return nil, err
	}

	return bytes.Bytes(), nil
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
