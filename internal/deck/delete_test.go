package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/mochi"
)

func Test_LeafDecks(t *testing.T) {
	decks := []mochi.Deck{
		{ID: "ROOT_1", ParentID: ""},
		{ID: "ROOT_2", ParentID: ""},
		{ID: "ROOT_3", ParentID: ""},
		{ID: "DECK_1", ParentID: "ROOT_1"},
		{ID: "DECK_2", ParentID: "ROOT_2"},
		{ID: "DECK_4", ParentID: "DECK_3"},
		{ID: "DECK_3", ParentID: "DECK_2"},
	}
	want := []string{"ROOT_3", "DECK_1", "DECK_4"}

	got := LeafDecks(decks)

	assert.ElementsMatch(t, want, got)
}
