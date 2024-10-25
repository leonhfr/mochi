package deck

import (
	"testing"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/test"
	"github.com/leonhfr/mochi/mochi"
)

func Test_Lock_cleanDecks(t *testing.T) {
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
				Decks: []map[string]lock.Deck{
					{
						"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
						"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
					},
				},
				UpdateDeck: []test.LockfileUpdateDeckName{
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
