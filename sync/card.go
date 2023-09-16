package sync

import (
	"context"
	"runtime"

	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
)

const (
	minHandlers = 4
	maxHandlers = 8
)

type CardResult struct {
	Created  int
	Updated  int
	Archived int
}

type syncCardResult struct {
	result CardResult
}

func SynchronizeCards(ctx context.Context, parsers []parser.Parser, sources []string, lock *Lock, config Config, client Client, fs filesystem.Interface) (CardResult, error) {
	jobMap, err := newJobMap(parsers, sources, lock, config)
	if err != nil {
		return CardResult{}, err
	}

	handlers := numHandlers()
	scr := &syncCardResult{}
	err = processJobMap(ctx, jobMap, handlers, scr, client, fs)

	return scr.result, err
}

type cardRequest struct{}

func (r *cardRequest) do(_ context.Context, _ *syncCardResult, _ Client) error {
	return nil
}

func generateCardRequests(_ context.Context, _ *deckJob, _ Client, _ filesystem.Interface) ([]cardRequest, error) {
	return nil, nil
}

func numHandlers() int {
	switch num := 2 * runtime.NumCPU(); {
	case num < minHandlers:
		return maxHandlers
	case num > maxHandlers:
		return maxHandlers
	default:
		return num
	}
}
