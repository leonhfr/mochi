package sync

import (
	"context"

	"github.com/leonhfr/mochi/api"
)

type Client interface {
	ListTemplates(ctx context.Context) ([]api.Template, error)
}
