package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/mochi"
)

// Client is the interface that should be implemented to sync cards.
type Client interface {
	CreateCard(ctx context.Context, req mochi.CreateCardRequest) (mochi.Card, error)
	UpdateCard(ctx context.Context, id string, req mochi.UpdateCardRequest) (mochi.Card, error)
	DeleteCard(ctx context.Context, id string) error
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
	Execute(ctx context.Context, client Client, reader image.Reader, lf Lockfile) error
}
