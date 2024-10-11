package card

import (
	"io"
	"path/filepath"

	"github.com/leonhfr/mochi/internal/parser"
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
