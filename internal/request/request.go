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
	SetCard(deckID string, cardID string, filename string) error
	GetImageHash(deckID, cardID, path string) (string, bool)
	SetImageHash(deckID, cardID, path, hash string) error
}

// Request is the interface that should be implemented to execute a request.
type Request interface {
	fmt.Stringer
	Sync(ctx context.Context, client Client, reader Reader, lf Lockfile) error
}
