package sync

import (
	"context"
	"runtime"
	"sync"

	"github.com/leonhfr/mochi/api"
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
	mu     sync.Mutex
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

type cardRequest struct {
	id      string
	archive bool
	create  *api.CreateCardRequest
	update  *api.UpdateCardRequest
}

func (r *cardRequest) do(ctx context.Context, scr *syncCardResult, client Client) error {
	if r.create != nil {
		if _, err := client.CreateCard(ctx, *r.create); err != nil {
			return err
		}

		scr.mu.Lock()
		scr.result.Created++
		scr.mu.Unlock()
	}

	if r.update != nil {
		if _, err := client.UpdateCard(ctx, r.id, *r.update); err != nil {
			return err
		}

		scr.mu.Lock()
		if r.archive {
			scr.result.Archived++
		} else {
			scr.result.Updated++
		}
		scr.mu.Unlock()
	}

	return nil
}

func generateCardRequests(ctx context.Context, job *deckJob, client Client, fs filesystem.Interface) ([]cardRequest, error) {
	_, err := parseCards(job, fs)
	if err != nil {
		return nil, err
	}

	_, err = client.ListCardsInDeck(ctx, job.id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func parseCards(job *deckJob, fs filesystem.Interface) ([]parser.Card, error) {
	var cards []parser.Card
	for _, source := range job.sources {
		content, err := fs.Read(source)
		if err != nil {
			return nil, err
		}

		parsedCards, err := job.parser.Convert(content)
		if err != nil {
			return nil, err
		}

		cards = append(cards, parsedCards...)
	}
	return cards, nil
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
