package request

import (
	"context"
	"fmt"

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
	SetCard(deckID string, cardID string, filename string) error
}

// Request is the interface that should be implemented to execute a request.
type Request interface {
	fmt.Stringer
	Sync(ctx context.Context, client Client, lf Lockfile) error
}
