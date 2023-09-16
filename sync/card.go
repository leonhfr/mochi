package sync

import (
	"context"

	"github.com/leonhfr/mochi/filesystem"
)

type CardResult struct {
	Created  int
	Updated  int
	Archived int
}

func SynchronizeCards(_ context.Context, _ []string, _ *Lock, _ Config, _ Client, _ filesystem.Interface) (CardResult, error) {
	return CardResult{}, nil
}
