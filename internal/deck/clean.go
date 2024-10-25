package deck

import (
	"context"
	"slices"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// CleanClient is the interface the mochi client should implement to clean the decks.
type CleanClient interface {
	ListDecks(ctx context.Context) ([]mochi.Deck, error)
}

// CleanLockfile is the interface the lockfile should implement to clean the decks.
type CleanLockfile interface {
	Lock()
	Unlock()
	Decks() map[string]lock.Deck
	UpdateDeck(id, name string)
	DeleteDeck(id string)
}

// Clean removes any decks from the lockfile that are not present in mochi.
func Clean(ctx context.Context, client CleanClient, lf CleanLockfile) error {
	mochiDecks, err := client.ListDecks(ctx)
	if err != nil {
		return err
	}

	cleanDecks(lf, mochiDecks)

	return nil
}

func cleanDecks(lf CleanLockfile, mochiDecks []mochi.Deck) {
	for deckID, deck := range lf.Decks() {
		index := slices.IndexFunc(mochiDecks, func(mochiDeck mochi.Deck) bool {
			return mochiDeck.ID == deckID
		})

		if index < 0 {
			lf.DeleteDeck(deckID)
			continue
		}

		if mochiDecks[index].ParentID != deck.ParentID {
			lf.DeleteDeck(deckID)
			continue
		}

		if mochiDecks[index].Name != deck.Name {
			lf.UpdateDeck(deckID, mochiDecks[index].Name)
		}
	}
}
