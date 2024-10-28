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
	images := image.New(reader, r.card.Path(), r.card.Images())

	filtered := filteredAttachments(lf, r.deckID, r.cardID, images)
	req := mochi.UpdateCardRequest{
		Content:     images.Replace(r.card.Content()),
		Attachments: filtered.Attachments(),
		Pos:         r.card.Position(),
	}

	if _, err := client.UpdateCard(ctx, r.cardID, req); err != nil {
		return err
	}

	lf.Lock()
	defer lf.Unlock()

	if err := lf.SetCard(r.deckID, r.cardID, r.card.Filename(), images.HashMap()); err != nil {
		return err
	}

	return nil
}

func filteredAttachments(lf Lockfile, deckID, cardID string, images image.Images) image.Images {
	lf.Lock()
	defer lf.Unlock()

	paths := make([]string, len(images))
	for i, image := range images {
		paths[i] = image.Path
	}

	card, ok := lf.Card(deckID, cardID)
	if !ok {
		return images
	}

	filtered := image.Images{}
	for _, image := range images {
		if md5, ok := card.Images[image.Path]; !ok || md5 != image.Hash {
			filtered = append(filtered, image)
		}
	}
	return filtered
}

// String implements the fmt.Stringer interface.
func (r *updateCard) String() string {
	return fmt.Sprintf("update request for card ID %s (%s)", r.cardID, r.card.Filename())
}
