package card

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

func Test_upsertSyncRequests(t *testing.T) {
	filename := "lorem-ipsum.md"
	deckID := "DECK_ID"
	mochiCards := []mochi.Card{
		{
			ID:      "CARD_ID_1",
			Name:    "CARD_TO_UPDATE",
			Content: "OLD_CONTENT",
		},
		{
			ID:      "CARD_ID_2",
			Name:    "CARD_TO_ARCHIVE",
			Content: "CONTENT",
		},
		{
			ID:      "CARD_ID_3",
			Name:    "CARD_TO_KEEP",
			Content: "CONTENT",
		},
	}
	parserCards := []parser.Card{
		{
			Name:     "CARD_TO_UPDATE",
			Content:  "NEW_CONTENT",
			Filename: filename,
		},
		{
			Name:     "CARD_TO_CREATE",
			Content:  "CONTENT",
			Filename: filename,
		},
		{
			Name:     "CARD_TO_KEEP",
			Content:  "CONTENT",
			Filename: filename,
		},
	}

	want := []Request{
		&updateCardRequest{
			cardID: "CARD_ID_1",
			req:    mochi.UpdateCardRequest{Content: "NEW_CONTENT"},
		},
		&archiveCardRequest{
			cardID: "CARD_ID_2",
			req:    mochi.UpdateCardRequest{Archived: true},
		},
		&createCardRequest{
			filename: filename,
			req: mochi.CreateCardRequest{
				Content: "CONTENT",
				DeckID:  "DECK_ID",
				Fields: map[string]mochi.Field{
					"name": {ID: "name", Value: "CARD_TO_CREATE"},
				},
			},
		},
	}

	got := upsertSyncRequests(filename, deckID, mochiCards, parserCards)
	assert.Equal(t, want, got)
}
