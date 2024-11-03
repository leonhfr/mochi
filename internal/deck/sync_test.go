package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/mochi"
)

func Test_upsertSyncRequests(t *testing.T) {
	path := "/testdata/lorem-ipsum.md"
	deckID := "DECK_ID"
	mochiCards := []mochi.Card{
		{
			ID:      "CARD_ID_1",
			Name:    "CARD_TO_UPDATE",
			Content: "OLD_CONTENT",
			Fields:  map[string]mochi.Field{"name": {ID: "name", Value: "CARD_TO_UPDATE"}},
		},
		{
			ID:      "CARD_ID_2",
			Name:    "CARD_TO_ARCHIVE",
			Content: "CONTENT",
			Fields:  map[string]mochi.Field{"name": {ID: "name", Value: "CARD_TO_ARCHIVE"}},
		},
		{
			ID:      "CARD_ID_3",
			Name:    "CARD_TO_KEEP",
			Content: "CONTENT",
			Fields:  map[string]mochi.Field{"name": {ID: "name", Value: "CARD_TO_KEEP"}},
		},
	}
	parserCards := []parser.Card{
		{
			Content: "NEW_CONTENT",
			Fields:  map[string]string{"name": "CARD_TO_UPDATE"},
			Path:    path,
		},
		{
			Content: "CONTENT",
			Fields:  map[string]string{"name": "CARD_TO_CREATE"},
			Path:    path,
		},
		{
			Content: "CONTENT",
			Fields:  map[string]string{"name": "CARD_TO_KEEP"},
			Path:    path,
		},
	}

	want := []request.Request{
		request.UpdateCard(deckID, "CARD_ID_1", parser.Card{
			Content: "NEW_CONTENT",
			Fields:  map[string]string{"name": "CARD_TO_UPDATE"},
			Path:    path,
		}, nil),
		request.DeleteCard("CARD_ID_2"),
		request.CreateCard("DECK_ID", parser.Card{
			Content: "CONTENT",
			Fields:  map[string]string{"name": "CARD_TO_CREATE"},
			Path:    path,
		}),
	}

	got := upsertSyncRequests(deckID, mochiCards, parserCards)
	assert.Equal(t, want, got)
}
