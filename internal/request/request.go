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
	AddAttachment(ctx context.Context, cardID, filename string, data []byte) error
}

// Lockfile is the interface the lockfile implement to sync cards.
type Lockfile interface {
	Lock()
	Unlock()
	SetCard(deckID, cardID, filename string) error
}

// Request is the interface that should be implemented to execute a request.
type Request interface {
	fmt.Stringer
	Execute(ctx context.Context, client Client, reader image.Reader, lf Lockfile) error
}

func mochiFields(fields map[string]string) map[string]mochi.Field {
	mochiFields := map[string]mochi.Field{}
	for key, value := range fields {
		mochiFields[key] = mochi.Field{ID: key, Value: value}
	}
	return mochiFields
}
