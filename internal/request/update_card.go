package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/converter"
	"github.com/leonhfr/mochi/mochi"
)

type updateCard struct {
	deckID      string
	cardID      string
	filename    string
	req         mochi.UpdateCardRequest
	attachments []converter.Attachment
}

// UpdateCard returns a new update card request.
func UpdateCard(deckID, cardID string, card card.Card, attachments map[string]mochi.Attachment) Request {
	return &updateCard{
		deckID:   deckID,
		cardID:   cardID,
		filename: card.Filename(),
		req: mochi.UpdateCardRequest{
			Content:    card.Content,
			TemplateID: card.TemplateID,
			Fields:     mochiFields(card.Fields),
			Pos:        card.Position,
		},
		attachments: filterAttachments(card.Attachments, attachments),
	}
}

// Execute implements the Request interface.
func (r *updateCard) Execute(ctx context.Context, client Client, lf Lockfile) error {
	if _, err := client.UpdateCard(ctx, r.cardID, r.req); err != nil {
		return err
	}

	for _, attachment := range r.attachments {
		if err := client.AddAttachment(ctx, r.cardID, attachment.Filename, attachment.Bytes); err != nil {
			return err
		}
	}

	lf.Lock()
	defer lf.Unlock()

	if err := lf.SetCard(r.deckID, r.cardID, r.filename); err != nil {
		return err
	}

	return nil
}

func filterAttachments(images []converter.Attachment, mochiAttachments map[string]mochi.Attachment) []converter.Attachment {
	attachments := []converter.Attachment{}
	for _, image := range images {
		if mochiAttachment, ok := mochiAttachments[image.Filename]; !ok || mochiAttachment.Size != len(image.Bytes) {
			attachments = append(attachments, image)
		}
	}
	return attachments
}

// String implements the fmt.Stringer interface.
func (r *updateCard) String() string {
	if len(r.attachments) > 0 {
		return fmt.Sprintf("update request for card ID %s (%s) with %d attachments", r.cardID, r.filename, len(r.attachments))
	}
	return fmt.Sprintf("update request for card ID %s (%s)", r.cardID, r.filename)
}
