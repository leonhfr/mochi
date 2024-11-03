package request

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type updateCard struct {
	deckID      string
	cardID      string
	card        parser.Card
	attachments map[string]mochi.Attachment
}

// UpdateCard returns a new update card request.
func UpdateCard(deckID, cardID string, card parser.Card, attachments map[string]mochi.Attachment) Request {
	return &updateCard{
		deckID:      deckID,
		cardID:      cardID,
		card:        card,
		attachments: attachments,
	}
}

// Execute implements the Request interface.
func (r *updateCard) Execute(ctx context.Context, client Client, reader image.Reader, lf Lockfile) error {
	images := image.New(reader, r.card.Path, r.card.Images)

	req := mochi.UpdateCardRequest{
		Content:    images.Replace(r.card.Content),
		TemplateID: r.card.TemplateID,
		Fields:     mochiFields(r.card.Fields),
		Pos:        r.card.Position,
	}

	if _, err := client.UpdateCard(ctx, r.cardID, req); err != nil {
		return err
	}

	p := pool.New().WithContext(ctx)
	for _, image := range filteredAttachments(images, r.attachments) {
		image := image
		p.Go(func(ctx context.Context) error {
			return client.AddAttachment(ctx, r.cardID, image.Filename, image.Bytes)
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

func filteredAttachments(images image.Images, attachments map[string]mochi.Attachment) image.Images {
	filtered := image.Images{}
	for _, image := range images {
		if attachment, ok := attachments[image.Filename]; !ok || attachment.Size != len(image.Bytes) {
			filtered = append(filtered, image)
		}
	}
	return filtered
}

// String implements the fmt.Stringer interface.
func (r *updateCard) String() string {
	return fmt.Sprintf("update request for card ID %s (%s)", r.cardID, r.card.Filename())
}
