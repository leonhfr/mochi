package request

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/converter"
	"github.com/leonhfr/mochi/mochi"
)

type updateCard struct {
	deckID      string
	cardID      string
	card        card.Card
	attachments map[string]mochi.Attachment
}

// UpdateCard returns a new update card request.
func UpdateCard(deckID, cardID string, card card.Card, attachments map[string]mochi.Attachment) Request {
	return &updateCard{
		deckID:      deckID,
		cardID:      cardID,
		card:        card,
		attachments: attachments,
	}
}

// Execute implements the Request interface.
func (r *updateCard) Execute(ctx context.Context, client Client, lf Lockfile) error {
	req := mochi.UpdateCardRequest{
		Content:    r.card.Content,
		TemplateID: r.card.TemplateID,
		Fields:     mochiFields(r.card.Fields),
		Pos:        r.card.Position,
	}

	if _, err := client.UpdateCard(ctx, r.cardID, req); err != nil {
		return err
	}

	p := pool.New().WithContext(ctx)
	for _, attachment := range filterAttachments(r.card.Attachments, r.attachments) {
		attachment := attachment
		p.Go(func(ctx context.Context) error {
			return client.AddAttachment(ctx, r.cardID, attachment.Filename, attachment.Bytes)
		})
	}

	err := p.Wait()
	if err != nil {
		return err
	}

	lf.Lock()
	defer lf.Unlock()

	if err := lf.SetCard(r.deckID, r.cardID, r.card.Filename()); err != nil {
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
	return fmt.Sprintf("update request for card ID %s (%s)", r.cardID, r.card.Filename())
}
