package request

import (
	"context"
	"fmt"
	"strings"

	"github.com/leonhfr/mochi/internal/image"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type create struct {
	deckID string
	card   parser.Card
}

// NewCreate returns a new create card request.
func NewCreate(deckID string, card parser.Card) Request {
	return &create{
		deckID: deckID,
		card:   card,
	}
}

// Sync implements the SyncRequest interface.
func (r *create) Sync(ctx context.Context, client Client, reader Reader, lf Lockfile) error {
	attachments, err := image.Attachments(reader, r.card.Images)
	if err != nil {
		return err
	}

	req := mochi.CreateCardRequest{
		Content: r.card.Content,
		DeckID:  r.deckID,
		Fields: map[string]mochi.Field{
			"name": {ID: "name", Value: r.card.Name},
		},
		Attachments: getAttachments(attachments),
		Pos:         getCardPos(r.card),
	}

	card, err := client.CreateCard(ctx, req)
	if err != nil {
		return err
	}

	hashMap := getHashMap(attachments)
	if err := lf.SetCard(r.deckID, card.ID, r.card.Filename, hashMap); err != nil {
		return err
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (r *create) String() string {
	return fmt.Sprintf("create request for file %s", r.card.Filename)
}

func getCardPos(card parser.Card) string {
	sanitized := strings.TrimFunc(card.Filename, func(r rune) bool {
		return ('0' <= r && r <= '9') || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
	})
	return fmt.Sprintf("%s%d", sanitized, card.Index)
}
