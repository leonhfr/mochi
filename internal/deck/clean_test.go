package deck

import (
	"testing"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/test"
	"github.com/leonhfr/mochi/mochi"
)

func Test_cleanDecks(t *testing.T) {
	tests := []struct {
		name  string
		decks []mochi.Deck
		calls test.Lockfile
	}{
		{
			name: "should not modify the lock",
			decks: []mochi.Deck{
				{ID: "DECK_ID_2", Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			calls: test.Lockfile{
				Lock: 1,
				Decks: []map[string]lock.Deck{
					{
						"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
						"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
					},
				},
			},
		},
		{
			name: "should remove decks that are not in the slice",
			decks: []mochi.Deck{
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			calls: test.Lockfile{
				Lock: 1,
				Decks: []map[string]lock.Deck{
					{
						"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
						"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
					},
				},
				DeleteDeck: []string{"DECK_ID_2"},
			},
		},
		{
			name: "should remove decks whose parent id have changed",
			decks: []mochi.Deck{
				{ID: "DECK_ID_2", Name: "DECK_NAME_2", ParentID: ""},
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			calls: test.Lockfile{
				Lock: 1,
				Decks: []map[string]lock.Deck{
					{
						"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
						"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
					},
				},
				DeleteDeck: []string{"DECK_ID_2"},
			},
		},
		{
			name: "should update the deck name",
			decks: []mochi.Deck{
				{ID: "DECK_ID_2", Name: "NEW_DECK_NAME_2", ParentID: "DECK_ID_1"},
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			calls: test.Lockfile{
				Lock: 1,
				Decks: []map[string]lock.Deck{
					{
						"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
						"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
					},
				},
				UpdateDeck: []test.LockfileUpdateDeck{
					{ID: "DECK_ID_2", Name: "NEW_DECK_NAME_2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lf := test.NewMockLockfile(tt.calls)
			cleanDecks(lf, tt.decks)
			lf.AssertExpectations(t)
		})
	}
}

func Test_cleanCards(t *testing.T) {
	tests := []struct {
		name   string
		deckID string
		cards  []mochi.Card
		calls  test.Lockfile
	}{
		{
			name:   "should not modify the lock if deck is not found",
			deckID: "DECK_ID_1",
			cards:  []mochi.Card{{ID: "CARD_ID_1"}, {ID: "CARD_ID_2"}, {ID: "CARD_ID_3"}},
			calls: test.Lockfile{
				Lock: 1,
				Deck: []test.LockfileDeck{
					{DeckID: "DECK_ID_1"},
				},
			},
		},
		{
			name:   "should not modify the lock if cards are present",
			deckID: "DECK_ID",
			cards:  []mochi.Card{{ID: "CARD_ID_1"}, {ID: "CARD_ID_2"}, {ID: "CARD_ID_3"}},
			calls: test.Lockfile{
				Lock: 1,
				Deck: []test.LockfileDeck{
					{DeckID: "DECK_ID", Deck: lock.Deck{Cards: map[string]lock.Card{
						"CARD_ID_1": {},
						"CARD_ID_2": {},
						"CARD_ID_3": {},
					}}, OK: true},
				},
			},
		},
		{
			name:   "should remove cards that are not in the slice",
			deckID: "DECK_ID",
			cards:  []mochi.Card{{ID: "CARD_ID_1"}, {ID: "CARD_ID_2"}},
			calls: test.Lockfile{
				Lock: 1,
				Deck: []test.LockfileDeck{
					{DeckID: "DECK_ID", Deck: lock.Deck{Cards: map[string]lock.Card{
						"CARD_ID_1": {},
						"CARD_ID_2": {},
						"CARD_ID_3": {},
					}}, OK: true},
				},
				DeleteCard: []test.LockfileDeleteCard{
					{DeckID: "DECK_ID", CardID: "CARD_ID_3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lf := test.NewMockLockfile(tt.calls)
			cleanCards(lf, tt.cards, tt.deckID)
			lf.AssertExpectations(t)
		})
	}
}
