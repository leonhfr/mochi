package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/internal/test"
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
		test.NewCard(test.ParserCard{
			Name:     "CARD_TO_UPDATE",
			Content:  "NEW_CONTENT",
			Filename: filename,
		}),
		test.NewCard(test.ParserCard{
			Name:     "CARD_TO_CREATE",
			Content:  "CONTENT",
			Filename: filename,
		}),
		test.NewCard(test.ParserCard{
			Name:     "CARD_TO_KEEP",
			Content:  "CONTENT",
			Filename: filename,
			Equals:   true,
		}),
	}

	want := []request.Request{
		request.UpdateCard(deckID, "CARD_ID_1", test.NewCard(test.ParserCard{
			Name:     "CARD_TO_UPDATE",
			Content:  "NEW_CONTENT",
			Filename: filename,
		}), nil),
		request.DeleteCard("CARD_ID_2"),
		request.CreateCard("DECK_ID", test.NewCard(test.ParserCard{
			Name:     "CARD_TO_CREATE",
			Content:  "CONTENT",
			Filename: filename,
		})),
	}

	got := upsertSyncRequests(deckID, mochiCards, parserCards)
	assert.Equal(t, want, got)
}
