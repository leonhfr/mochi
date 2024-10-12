package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/parser/image"
	"github.com/leonhfr/mochi/mochi"
)

type update struct {
	deckID string
	cardID string
	card   parser.Card
}

// NewUpdate returns a new update card request.
func NewUpdate(deckID, cardID string, card parser.Card) Request {
	return &update{
		deckID: deckID,
		cardID: cardID,
		card:   card,
	}
}

// Sync implements the SyncRequest interface.
func (r *update) Sync(ctx context.Context, client Client, reader Reader, lf Lockfile) error {
	attachments, err := r.card.Images.Attachments(reader)
	if err != nil {
		return err
	}

	lf.CleanImages(r.deckID, r.cardID, getPaths(attachments))

	filtered := []image.Attachment{}
	for _, attachment := range attachments {
		hash, ok := lf.GetImageHash(r.deckID, r.cardID, attachment.Path)
		if !ok || hash != attachment.Hash {
			filtered = append(filtered, attachment)
		}
	}

	req := mochi.UpdateCardRequest{
		Content:     r.card.Content,
		Attachments: getAttachments(filtered),
	}

	if _, err := client.UpdateCard(ctx, r.cardID, req); err != nil {
		return err
	}

	for _, attachment := range attachments {
		if err := lf.SetImageHash(r.deckID, r.cardID, attachment.Path, attachment.Hash); err != nil {
			return err
		}
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *update) String() string {
	return fmt.Sprintf("update request for card ID %s", r.cardID)
}
