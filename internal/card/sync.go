package card

import (
	"slices"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/request"
	"github.com/leonhfr/mochi/mochi"
)

// Lockfile is the interface that should be implemented to update the lockfile.
type Lockfile interface {
	GetCard(deckID string, cardID string) (lock.Card, bool)
}

// SyncRequests parses the note files and returns the requests
// required to sync them.
func SyncRequests(lf Lockfile, deckID string, mochiCards []mochi.Card, parsedCards []parser.Card) []request.Request {
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
		index := slices.IndexFunc(tmp, indexFunc(mochiCard))
		if index < 0 {
			reqs = append(reqs, request.DeleteCard(mochiCard.ID))
			continue
		}

		if mochiCard.Content != tmp[index].Content {
			reqs = append(reqs, request.UpdateCard(deckID, mochiCard.ID, tmp[index]))
		}
		tmp = sliceRemove(tmp, index)
	}

	for _, parsedCard := range tmp {
		reqs = append(reqs, request.CreateCard(deckID, parsedCard))
	}

	return reqs
}

func indexFunc(mochiCard mochi.Card) func(c parser.Card) bool {
	return func(parsedCard parser.Card) bool {
		name, ok := mochiCard.Fields["name"]
		if !ok {
			return mochiCard.Name == parsedCard.Name
		}
		return name.Value == parsedCard.Name
	}
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

func groupMochiCardsByFilename(lf Lockfile, deckID string, mochiCards []mochi.Card) (map[string][]mochi.Card, []mochi.Card) {
	matched := make(map[string][]mochi.Card)
	var notMatched []mochi.Card
	for _, mochiCard := range mochiCards {
		if lockCard, ok := lf.GetCard(deckID, mochiCard.ID); ok {
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
		matched[parsedCard.Filename] = append(matched[parsedCard.Filename], parsedCard)
	}
	return matched
}
