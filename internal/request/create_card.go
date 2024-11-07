package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/converter"
	"github.com/leonhfr/mochi/mochi"
)

type createRequest struct {
	deckID      string
	filename    string
	req         mochi.CreateCardRequest
	attachments []converter.Attachment
}

// CreateCard returns a new create card request.
func CreateCard(deckID string, card card.Card) Request {
	return &createRequest{
		deckID:   deckID,
		filename: card.Filename(),
		req: mochi.CreateCardRequest{
			Content:    card.Content,
			DeckID:     deckID,
			TemplateID: card.TemplateID,
			Fields:     mochiFields(card.Fields),
			Pos:        card.Position,
		},
		attachments: card.Attachments,
	}
}

// Execute implements the Request interface.
func (r *createRequest) Execute(ctx context.Context, client Client, lf Lockfile) error {
	card, err := client.CreateCard(ctx, r.req)
	if err != nil {
		return err
	}

	for _, attachment := range r.attachments {
		if err := client.AddAttachment(ctx, card.ID, attachment.Filename, attachment.Bytes); err != nil {
			return err
		}
	}

	lf.Lock()
	defer lf.Unlock()

	if err := lf.SetCard(r.deckID, card.ID, r.filename); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *createRequest) String() string {
	if len(r.attachments) > 0 {
		return fmt.Sprintf("create request for file %s with %d attachments", r.filename, len(r.attachments))
	}
	return fmt.Sprintf("create request for file %s", r.filename)
}
