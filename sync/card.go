package sync

import (
	"context"
	"runtime"

	"golang.org/x/exp/slices"

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

func SynchronizeCards(ctx context.Context, parsers []parser.Parser, sources []string, lock *Lock, config Config, client Client, fs filesystem.Interface) (CardResult, error) {
	jobMap, err := newJobMap(parsers, sources, lock, config)
	if err != nil {
		return CardResult{}, err
	}

	handlers := numHandlers()

	return processJobMap(ctx, jobMap, handlers, lock, client, fs)
}

func generateCardRequests(ctx context.Context, job *deckJob, lock *Lock, client Client, fs filesystem.Interface) ([]cardRequest, error) {
	cards, err := parseCards(job, fs)
	if err != nil {
		return nil, err
	}

	apiCards, err := client.ListCardsInDeck(ctx, job.id)
	if err != nil {
		return nil, err
	}

	var requests []cardRequest
	for _, card := range cards {
		index := slices.IndexFunc[[]api.Card](apiCards, func(apiCard api.Card) bool {
			return card.Name == apiCard.Name
		})

		if index < 0 {
			request, err := newCreateCardRequest(job, card, fs)
			if err != nil {
				return nil, err
			}
			requests = append(requests, request)
			continue
		}

		apiCard := apiCards[index]
		apiCards = append(apiCards[:index], apiCards[index+1:]...)
		if !cardEqual(job, card, apiCard) {
			request, err := newUpdateCardRequest(job, apiCard.ID, card, lock, fs)
			if err != nil {
				return nil, err
			}
			requests = append(requests, request)
		}
	}

	if job.archive {
		for _, apiCard := range apiCards {
			requests = append(requests, newArchiveCardRequest(apiCard.ID))
		}
	}

	return requests, nil
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

func cardEqual(job *deckJob, card parser.Card, apiCard api.Card) bool {
	if !job.hasTemplate {
		return card.Name == apiCard.Name && card.Content == apiCard.Content
	}

	if card.Name != apiCard.Name || job.template.TemplateID != apiCard.TemplateID {
		return false
	}

	if len(card.Fields) != len(apiCard.Fields) {
		return false
	}

	for id, field := range job.template.Fields {
		if card.Fields[field] != apiCard.Fields[id].Value {
			return false
		}
	}

	return true
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
