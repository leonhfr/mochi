package deck

import (
	"slices"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/mochi"
)

// SyncLockfile is the interface that should be implemented to update the lockfile.
type SyncLockfile interface {
	Lock()
	Unlock()
	Card(deckID string, cardID string) (lock.Card, bool)
}

// SyncRequests parses the note files and returns the requests
// required to sync them.
func SyncRequests(lf SyncLockfile, deckID string, mochiCards []mochi.Card, parsedCards []parser.Card) []request.Request {
	groupedMochiCards, notMatched := groupMochiCardsByFilename(lf, deckID, mochiCards)
	groupedParsedCards := groupParsedCardsByFilename(parsedCards)
	groupedCards := groupCardsByFilename(groupedMochiCards, groupedParsedCards)
	reqs := []request.Request{}
	for _, mochiCard := range notMatched {
		reqs = append(reqs, request.DeleteCard(mochiCard.ID))
	}
	for _, group := range groupedCards {
		groupReqs := upsertSyncRequests(deckID, group.mochi, group.parsed)
		reqs = append(reqs, groupReqs...)
	}
	return reqs
}

type fileGroup struct {
	mochi  []mochi.Card
	parsed []parser.Card
}

func upsertSyncRequests(deckID string, mochiCards []mochi.Card, parsedCards []parser.Card) []request.Request {
	tmp := make([]parser.Card, len(parsedCards))
	copy(tmp, parsedCards)

	reqs := []request.Request{}
	for _, mochiCard := range mochiCards {
		index := slices.IndexFunc(tmp, func(card parser.Card) bool { return card.Is(mochiCard) })
		if index < 0 {
			reqs = append(reqs, request.DeleteCard(mochiCard.ID))
			continue
		}

		if !tmp[index].Equals(mochiCard) {
			reqs = append(reqs, request.UpdateCard(deckID, mochiCard.ID, tmp[index]))
		}
		tmp = sliceRemove(tmp, index)
	}

	for _, parsedCard := range tmp {
		reqs = append(reqs, request.CreateCard(deckID, parsedCard))
	}

	return reqs
}

func sliceRemove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func groupCardsByFilename(mochiCards map[string][]mochi.Card, parsedCards map[string][]parser.Card) map[string]fileGroup {
	groups := make(map[string]fileGroup)
	for filename, cards := range mochiCards {
		groups[filename] = fileGroup{mochi: cards}
	}
	for filename, cards := range parsedCards {
		groups[filename] = fileGroup{
			mochi:  groups[filename].mochi,
			parsed: cards,
		}
	}
	return groups
}

func groupMochiCardsByFilename(lf SyncLockfile, deckID string, mochiCards []mochi.Card) (map[string][]mochi.Card, []mochi.Card) {
	lf.Lock()
	defer lf.Unlock()

	matched := make(map[string][]mochi.Card)
	var notMatched []mochi.Card
	for _, mochiCard := range mochiCards {
		if lockCard, ok := lf.Card(deckID, mochiCard.ID); ok {
			matched[lockCard.Filename] = append(matched[lockCard.Filename], mochiCard)
		} else {
			notMatched = append(notMatched, mochiCard)
		}
	}
	return matched, notMatched
}

func groupParsedCardsByFilename(parsedCards []parser.Card) map[string][]parser.Card {
	matched := make(map[string][]parser.Card)
	for _, parsedCard := range parsedCards {
		matched[parsedCard.Filename()] = append(matched[parsedCard.Filename()], parsedCard)
	}
	return matched
}
