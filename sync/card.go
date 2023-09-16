package sync

import (
	"context"

	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
)

type CardResult struct {
	Created  int
	Updated  int
	Archived int
}

func SynchronizeCards(_ context.Context, parsers []parser.Parser, sources []string, lock *Lock, config Config, _ Client, _ filesystem.Interface) (CardResult, error) {
	_, err := newJobMap(parsers, sources, lock, config)
	if err != nil {
		return CardResult{}, err
	}
	return CardResult{}, nil
}
