package card

import (
	"io"
	"path/filepath"
	"slices"

	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

// Reader represents the interface to read note files.
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// Parser represents the interface to parse note files.
type Parser interface {
	Convert(parser, path string, source io.Reader) ([]parser.Card, error)
}

// Parse parses the note files for cards.
func Parse(r Reader, p Parser, workspace, parserName string, filePaths []string) ([]parser.Card, error) {
	var cards []parser.Card
	for _, path := range filePaths {
		parsed, err := parseFile(r, p, workspace, parserName, path)
		if err != nil {
			return nil, err
		}
		cards = append(cards, parsed...)
	}
	return cards, nil
}

func parseFile(r Reader, p Parser, workspace, parserName, path string) ([]parser.Card, error) {
	path = filepath.Join(workspace, path)
	bytes, err := r.Read(path)
	if err != nil {
		return nil, err
	}
	defer bytes.Close()

	cards, err := p.Convert(parserName, path, bytes)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// ReadLockfile is the interface to interact with the lockfile.
type ReadLockfile interface {
	GetCard(deckID string, cardID string) (lock.Card, bool)
}

// SyncRequests parses the note files and returns the requests
// required to sync them.
func SyncRequests(lf ReadLockfile, deckID string, mochiCards []mochi.Card, parsedCards []parser.Card) []SyncRequest {
	groupedMochiCards, notMatched := groupMochiCardsByFilename(lf, deckID, mochiCards)
	groupedParsedCards := groupParsedCardsByFilename(parsedCards)
	groupedCards := groupCardsByFilename(groupedMochiCards, groupedParsedCards)
	reqs := []SyncRequest{}
	for _, mochiCard := range notMatched {
		reqs = append(reqs, newArchiveCardRequest(mochiCard.ID))
	}
	for filename, group := range groupedCards {
		createReqs, updateReqs, archiveReqs := upsertSyncRequests(filename, deckID, group.mochi, group.parsed)
		for _, r := range createReqs {
			reqs = append(reqs, r)
		}
		for _, r := range updateReqs {
			reqs = append(reqs, r)
		}
		for _, r := range archiveReqs {
			reqs = append(reqs, r)
		}
	}
	return reqs
}

type fileGroup struct {
	mochi  []mochi.Card
	parsed []parser.Card
}

func upsertSyncRequests(filename, deckID string, mochiCards []mochi.Card, parsedCards []parser.Card) ([]*createCardRequest, []*updateCardRequest, []*archiveCardRequest) {
	tmp := make([]parser.Card, len(parsedCards))
	copy(tmp, parsedCards)

	createReqs := []*createCardRequest{}
	updateReqs := []*updateCardRequest{}
	archiveReqs := []*archiveCardRequest{}

	for _, mochiCard := range mochiCards {
		index := slices.IndexFunc(tmp, indexFunc(mochiCard))
		if index < 0 {
			archiveReqs = append(archiveReqs, newArchiveCardRequest(mochiCard.ID))
			continue
		}

		if mochiCard.Content != tmp[index].Content {
			updateReqs = append(updateReqs, newUpdateCardRequest(mochiCard.ID, tmp[index].Content))
		}
		tmp = sliceRemove(tmp, index)
	}

	for _, parsedCard := range tmp {
		createReqs = append(createReqs, newCreateCardRequest(filename, deckID, parsedCard.Name, parsedCard.Content))
	}

	return createReqs, updateReqs, archiveReqs
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

func groupMochiCardsByFilename(lf ReadLockfile, deckID string, mochiCards []mochi.Card) (map[string][]mochi.Card, []mochi.Card) {
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
