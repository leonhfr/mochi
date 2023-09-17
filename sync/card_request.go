package sync

import (
	"context"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/parser"
)

type requestKind int

const (
	createRequest requestKind = iota
	updateRequest
	archiveRequest
)

func (cr *CardResult) increment(kind requestKind) {
	switch kind {
	case createRequest:
		cr.Created++
	case updateRequest:
		cr.Updated++
	case archiveRequest:
		cr.Archived++
	}
}

type cardRequest struct {
	kind       requestKind
	id         string
	deckID     string
	content    string
	templateID string
	archived   bool
	fields     map[string]api.Field
}

func newCreateCardRequest(job *deckJob, card parser.Card) cardRequest {
	content, fields := newCardContent(job, card)
	return cardRequest{
		kind:       createRequest,
		deckID:     job.id,
		content:    content,
		templateID: job.template.TemplateID,
		fields:     fields,
	}
}

func newUpdateCardRequest(job *deckJob, id string, card parser.Card) cardRequest {
	content, fields := newCardContent(job, card)
	return cardRequest{
		kind:       updateRequest,
		id:         id,
		deckID:     job.id,
		content:    content,
		templateID: job.template.TemplateID,
		fields:     fields,
	}
}

func newArchiveCardRequest(id string) cardRequest {
	return cardRequest{
		id:       id,
		kind:     archiveRequest,
		archived: true,
	}
}

func processCardRequest(ctx context.Context, req cardRequest, client Client) error {
	switch req.kind {
	case createRequest:
		_, err := client.CreateCard(ctx, api.CreateCardRequest{
			DeckID:     req.deckID,
			Content:    req.content,
			TemplateID: req.templateID,
			Fields:     req.fields,
		})
		return err
	case updateRequest:
		_, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
			DeckID:     req.deckID,
			Content:    req.content,
			TemplateID: req.templateID,
			Fields:     req.fields,
		})
		return err
	case archiveRequest:
		_, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
			Archived: true,
		})
		return err
	default:
		return nil
	}
}

func newCardContent(job *deckJob, card parser.Card) (string, map[string]api.Field) {
	if !job.hasTemplate {
		return card.Content, nil
	}

	fields := make(map[string]api.Field)
	for id, field := range job.template.Fields {
		fields[id] = api.Field{
			ID:    id,
			Value: card.Fields[field],
		}
	}
	return "", fields
}
