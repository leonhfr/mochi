package sync

import (
	"context"
	"runtime"
	"sync"

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

type requestKind int

const (
	createRequest requestKind = iota
	updateRequest
	archiveRequest
)

type cardRequest struct {
	id     string
	kind   requestKind
	create api.CreateCardRequest
	update api.UpdateCardRequest
}

func (r *cardRequest) do(ctx context.Context, scr *syncCardResult, client Client) error {
	if r.kind == createRequest {
		if _, err := client.CreateCard(ctx, r.create); err != nil {
			return err
		}

		scr.mu.Lock()
		scr.result.Created++
		scr.mu.Unlock()
	}

	if r.kind == updateRequest || r.kind == archiveRequest {
		if _, err := client.UpdateCard(ctx, r.id, r.update); err != nil {
			return err
		}

		scr.mu.Lock()
		if r.kind == updateRequest {
			scr.result.Updated++
		} else {
			scr.result.Archived++
		}
		scr.mu.Unlock()
	}

	return nil
}

func generateCardRequests(ctx context.Context, job *deckJob, client Client, fs filesystem.Interface) ([]cardRequest, error) {
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
			requests = append(requests, newCreateCardRequest(job, card))
			continue
		}

		apiCard := apiCards[index]
		apiCards = append(apiCards[:index], apiCards[index+1:]...)
		if !cardEqual(job, card, apiCard) {
			requests = append(requests, newUpdateCardRequest(job, apiCard.ID, card))
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

func newCreateCardRequest(deck *deckJob, card parser.Card) cardRequest {
	if !deck.hasTemplate {
		return cardRequest{
			kind: createRequest,
			create: api.CreateCardRequest{
				Content: card.Content,
				DeckID:  deck.id,
			},
		}
	}

	fields := make(map[string]api.Field)
	for id, field := range deck.template.Fields {
		fields[id] = api.Field{
			ID:    id,
			Value: card.Fields[field],
		}
	}

	return cardRequest{
		kind: createRequest,
		create: api.CreateCardRequest{
			Fields:     fields,
			DeckID:     deck.id,
			TemplateID: deck.template.TemplateID,
		},
	}
}

func newUpdateCardRequest(job *deckJob, id string, card parser.Card) cardRequest {
	if !job.hasTemplate {
		return cardRequest{
			id:   id,
			kind: updateRequest,
			update: api.UpdateCardRequest{
				Content: card.Content,
				DeckID:  job.id,
			},
		}
	}

	fields := make(map[string]api.Field)
	for id, field := range job.template.Fields {
		fields[id] = api.Field{
			ID:    id,
			Value: card.Fields[field],
		}
	}

	return cardRequest{
		id:   id,
		kind: updateRequest,
		update: api.UpdateCardRequest{
			Fields:     fields,
			DeckID:     job.id,
			TemplateID: job.template.TemplateID,
		},
	}
}

func newArchiveCardRequest(id string) cardRequest {
	return cardRequest{
		id:     id,
		kind:   archiveRequest,
		update: api.UpdateCardRequest{Archived: true},
	}
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
