package sync

import (
	"context"

	"github.com/leonhfr/mochi/api"
)

type Client interface {
	ListDecks(ctx context.Context) ([]api.Deck, error)
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
