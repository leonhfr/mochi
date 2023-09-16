package sync

import (
	"context"

	"github.com/leonhfr/mochi/api"
)

type Client interface {
	ListDecks(ctx context.Context) ([]api.Deck, error)
	CreateDeck(ctx context.Context, req api.CreateDeckRequest) (api.Deck, error)
	UpdateDeck(ctx context.Context, id string, req api.UpdateDeckRequest) (api.Deck, error)
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
