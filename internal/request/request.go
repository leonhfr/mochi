package request

import (
	"context"
	"fmt"
	"io"

	"github.com/leonhfr/mochi/mochi"
)

// Client is the interface that should be implemented to sync cards.
type Client interface {
	CreateCard(ctx context.Context, req mochi.CreateCardRequest) (mochi.Card, error)
	UpdateCard(ctx context.Context, id string, req mochi.UpdateCardRequest) (mochi.Card, error)
	DeleteCard(ctx context.Context, id string) error
}

// Reader represents the interface to read files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Lockfile is the interface the lockfile implement to sync cards.
type Lockfile interface {
	CleanImages(deckID, cardID string, paths []string)
	SetCard(deckID, cardID, filename string, images map[string]string) error
	GetImageHashes(deckID, cardID string, paths []string) []string
}

// Request is the interface that should be implemented to execute a request.
type Request interface {
	fmt.Stringer
	Sync(ctx context.Context, client Client, reader Reader, lf Lockfile) error
}

func getAttachments(images []mochiAttachment) []mochi.Attachment {
	attachments := make([]mochi.Attachment, 0, len(images))
	for _, image := range images {
		attachments = append(attachments, image.Mochi)
	}
	return attachments
}

func getPaths(images []mochiAttachment) []string {
	paths := make([]string, 0, len(images))
	for _, image := range images {
		paths = append(paths, image.Path)
	}
	return paths
}

func getHashMap(images []mochiAttachment) map[string]string {
	hashMap := make(map[string]string)
	for _, image := range images {
		hashMap[image.Path] = image.Hash
	}
	return hashMap
}
