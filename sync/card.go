package sync

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"

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

func SynchronizeCards(ctx context.Context, parsers []parser.Parser, sources []string, lock *Lock, config Config, client Client, fs filesystem.Interface, logger Logger) (CardResult, error) {
	jobMap, err := newJobMap(parsers, sources, lock, config)
	if err != nil {
		return CardResult{}, err
	}

	var logs []string
	for path, job := range jobMap {
		if job.hasTemplate {
			logs = append(logs, fmt.Sprintf("%s: %d sources to sync (template id: %s)", path, len(job.sources), job.template.TemplateID))
		} else {
			logs = append(logs, fmt.Sprintf("%s: %d sources to sync (parser: %s)", path, len(job.sources), job.parser.String()))
		}
	}
	sort.Strings(logs)
	for _, log := range logs {
		logger.Info(log)
	}

	handlers := numHandlers()

	return processJobMap(ctx, jobMap, handlers, lock, client, fs, logger)
}

func generateCardRequests(ctx context.Context, job *deckJob, lock *Lock, client Client, fs filesystem.Interface) ([]cardRequest, error) {
	cardsMap, err := parseCards(job, fs)
	if err != nil {
		return nil, err
	}

	apiCards, err := client.ListCardsInDeck(ctx, job.id)
	if err != nil {
		return nil, err
	}

	cardIDs := make([]string, 0, len(apiCards))
	apiCardsMap := make(map[string]api.Card)
	for _, apiCard := range apiCards {
		cardIDs = append(cardIDs, apiCard.ID)
		apiCardsMap[apiCard.ID] = apiCard
	}
	lock.cleanCards(job.id, cardIDs)

	var requests []cardRequest
	for filename, cards := range cardsMap {
		for _, card := range cards {
			id, reqs, err := matchCard(filename, card, apiCardsMap, job, lock, fs)
			if err != nil {
				return nil, err
			}

			if len(reqs) > 0 {
				requests = append(requests, reqs...)
			}
			if len(id) > 0 {
				delete(apiCardsMap, id)
			}
		}
	}

	if job.archive {
		for id := range apiCardsMap {
			requests = append(requests, newArchiveCardRequest(id))
		}
	}

	return requests, nil
}

func matchCard(filename string, card parser.Card, apiCardsMap map[string]api.Card, job *deckJob, lock *Lock, fs filesystem.Interface) (string, []cardRequest, error) {
	matchedCards := make(map[string]api.Card)
	for _, apiCard := range apiCardsMap {
		if card.Name == apiCard.Name {
			matchedCards[apiCard.ID] = apiCard
		}
	}

	for id, apiCard := range matchedCards {
		if lockCard, ok := lock.getCard(job.id, id); ok && lockCard.Filename == filename {
			if cardEqual(job, card, apiCard) {
				return id, nil, nil
			}

			request, err := newUpdateCardRequest(job, apiCard.ID, filename, card, lock, fs)
			if err != nil {
				return "", nil, err
			}

			return id, []cardRequest{request}, nil
		}
	}

	request, err := newCreateCardRequest(job, filename, card, fs)
	if err != nil {
		return "", nil, err
	}

	return "", []cardRequest{request}, nil
}

func parseCards(job *deckJob, fs filesystem.Interface) (map[string][]parser.Card, error) {
	cards := make(map[string][]parser.Card)
	for _, source := range job.sources {
		content, err := fs.Read(source)
		if err != nil {
			return nil, err
		}

		filename := filepath.Base(source)
		parsedCards, err := job.parser.Convert(source, content)
		if err != nil {
			return nil, err
		}

		cards[filename] = parsedCards
	}
	return cards, nil
}

func cardEqual(job *deckJob, card parser.Card, apiCard api.Card) bool {
	if !job.hasTemplate {
		return card.Name == apiCard.Name && card.Content == apiCard.Content && !apiCard.Archived
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

	return !apiCard.Archived
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
