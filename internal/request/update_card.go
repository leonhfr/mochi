package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type updateCard struct {
	deckID string
	cardID string
	card   parser.Card
}

// UpdateCard returns a new update card request.
func UpdateCard(deckID, cardID string, card parser.Card) Request {
	return &updateCard{
		deckID: deckID,
		cardID: cardID,
		card:   card,
	}
}

// Execute implements the Request interface.
func (r *updateCard) Execute(ctx context.Context, client Client, reader image.Reader, lf Lockfile) error {
	images := image.New(reader, r.card.Path, r.card.Images)

	paths := images.Paths()
	lf.CleanImages(r.deckID, r.cardID, paths)

	hashes := lf.GetImageHashes(r.deckID, r.cardID, paths)
	filtered := image.Images{}
	for index, image := range images {
		if hash := hashes[index]; hash != "" || hash != image.Hash {
			filtered = append(filtered, image)
		}
	}

	req := mochi.UpdateCardRequest{
		Content:     images.Replace(r.card.Content),
		Attachments: filtered.Attachments(),
	}

	if _, err := client.UpdateCard(ctx, r.cardID, req); err != nil {
		return err
	}

	hashMap := filtered.HashMap()
	if err := lf.SetCard(r.deckID, r.cardID, r.card.Filename, hashMap); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *updateCard) String() string {
	return fmt.Sprintf("update request for card ID %s", r.cardID)
}
