package image

import (
	"bytes"
	//nolint:gosec
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/leonhfr/mochi/mochi"
)

// Reader represents the interface to read files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Attachment contains the data of an image attachment.
type Attachment struct {
	Mochi mochi.Attachment
	Hash  string
	Path  string
}

// Attachments returns the images attachments.
func Attachments(r Reader, images map[string]Image) ([]Attachment, error) {
	attachments := make([]Attachment, 0, len(images))
	for path, image := range images {
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
		FileName:    fmt.Sprintf("%s.%s", image.Filename, image.Extension),
		ContentType: image.MimeType,
		Data:        bytes.String(),
	}, nil
}
