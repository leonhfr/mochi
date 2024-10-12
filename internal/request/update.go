package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/internal/parser"
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
	attachments, err := image.Attachments(reader, r.card.Images)
	if err != nil {
		return err
	}

	paths := getPaths(attachments)
	lf.CleanImages(r.deckID, r.cardID, paths)

	hashes := lf.GetImageHashes(r.deckID, r.cardID, paths)
	filtered := []image.Attachment{}
	for index, attachment := range attachments {
		if hash := hashes[index]; hash != "" || hash != attachment.Hash {
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

	hashMap := getHashMap(attachments)
	if err := lf.SetCard(r.deckID, r.cardID, r.card.Filename, hashMap); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *update) String() string {
	return fmt.Sprintf("update request for card ID %s", r.cardID)
}
