package deck

import (
	"context"
	"slices"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/mochi"
)

// CleanDecksClient is the interface the mochi client should implement to clean the decks.
type CleanDecksClient interface {
	ListDecks(ctx context.Context) ([]mochi.Deck, error)
}

// CleanCardsClient is the interface the mochi client should implement to generate the sync requests.
type CleanCardsClient interface {
	ListCardsInDeck(ctx context.Context, deckID string) ([]mochi.Card, error)
}

// CleanDecksLockfile is the interface the lockfile should implement to clean the decks.
type CleanDecksLockfile interface {
	Lock()
	Unlock()
	Decks() map[string]lock.Deck
	UpdateDeck(id, name string)
	DeleteDeck(id string)
}

// CleanCardsLockfile is the interface the lockfile should implement to clean the cards.
type CleanCardsLockfile interface {
	Lock()
	Unlock()
	Deck(id string) (lock.Deck, bool)
	// Card(deckID string, cardID string) (lock.Card, bool)
	DeleteCard(deckID, cardID string)
}

// CleanDecks removes any decks from the lockfile that are not present in mochi.
func CleanDecks(ctx context.Context, client CleanDecksClient, lf CleanDecksLockfile) error {
	mochiDecks, err := client.ListDecks(ctx)
	if err != nil {
		return err
	}

	cleanDecks(lf, mochiDecks)
	return nil
}

// CleanCards removes any cards from the lockfile that are not present in mochi.
func CleanCards(ctx context.Context, client CleanCardsClient, lf CleanCardsLockfile, deckID string) error {
	mochiCards, err := client.ListCardsInDeck(ctx, deckID)
	if err != nil {
		return err
	}

	cleanCards(lf, mochiCards, deckID)
	return nil
}

func cleanDecks(lf CleanDecksLockfile, mochiDecks []mochi.Deck) {
	lf.Lock()
	defer lf.Unlock()

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

func cleanCards(lf CleanCardsLockfile, mochiCards []mochi.Card, deckID string) {
	lf.Lock()
	defer lf.Unlock()

	deck, ok := lf.Deck(deckID)
	if !ok {
		return
	}

	for cardID := range deck.Cards {
		index := slices.IndexFunc(mochiCards, func(mochiCard mochi.Card) bool {
			return mochiCard.ID == cardID
		})

		if index < 0 {
			lf.DeleteCard(deckID, cardID)
		}
	}
}
