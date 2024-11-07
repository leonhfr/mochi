package deck

import (
	"slices"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/converter"
	"github.com/leonhfr/mochi/internal/lock"
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
func SyncRequests(lf SyncLockfile, deckID string, mochiCards []mochi.Card, parsedCards []card.Card) []request.Request {
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
	parsed []card.Card
}

func upsertSyncRequests(deckID string, mochiCards []mochi.Card, parsedCards []card.Card) []request.Request {
	tmp := make([]card.Card, len(parsedCards))
	copy(tmp, parsedCards)

	reqs := []request.Request{}
	for _, mochiCard := range mochiCards {
		index := slices.IndexFunc(tmp, func(card card.Card) bool { return cardIs(card, mochiCard) })
		if index < 0 {
			reqs = append(reqs, request.DeleteCard(mochiCard.ID))
			continue
		}

		if !cardEquals(tmp[index], mochiCard) {
			reqs = append(reqs, request.UpdateCard(deckID, mochiCard.ID, tmp[index], mochiCard.Attachments))
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

func groupCardsByFilename(mochiCards map[string][]mochi.Card, parsedCards map[string][]card.Card) map[string]fileGroup {
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

func groupParsedCardsByFilename(parsedCards []card.Card) map[string][]card.Card {
	matched := make(map[string][]card.Card)
	for _, parsedCard := range parsedCards {
		matched[parsedCard.Filename()] = append(matched[parsedCard.Filename()], parsedCard)
	}
	return matched
}

func cardIs(card card.Card, mochiCard mochi.Card) bool {
	name, ok := mochiCard.Fields["name"]
	return ok && name.Value == card.Fields["name"]
}

func cardEquals(card card.Card, mochiCard mochi.Card) bool {
	return mochiCard.Content == card.Content &&
		mochiCard.TemplateID == card.TemplateID &&
		mochiCard.Pos == card.Position &&
		mapsEqual(mochiCard.Fields, mochiFields(card.Fields)) &&
		hasAttachments(card.Attachments, mochiCard.Attachments)
}

func mochiFields(fields map[string]string) map[string]mochi.Field {
	mochiFields := map[string]mochi.Field{}
	for key, value := range fields {
		mochiFields[key] = mochi.Field{ID: key, Value: value}
	}
	return mochiFields
}

func hasAttachments(attachments []converter.Attachment, mochiAttachments map[string]mochi.Attachment) bool {
	for _, attachment := range attachments {
		if mochiAttachment, ok := mochiAttachments[attachment.Filename]; !ok || mochiAttachment.Size != len(attachment.Bytes) {
			return false
		}
	}
	return true
}

func mapsEqual[T comparable](m1, m2 map[string]T) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || v1 != v2 {
			return false
		}
	}
	return true
}
