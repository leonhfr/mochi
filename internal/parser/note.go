package parser

import "github.com/leonhfr/mochi/internal/parser/image"

// note represents a note parser.
//
// The whole content of the file is returned as a card.
// The is the file name without the extension.
type note struct {
	fc FileCheck
}

// newNote returns a new note parser.
func newNote(fc FileCheck) *note {
	return &note{fc: fc}
}

// Convert implements the cardParser interface.
func (n *note) convert(path string, source []byte) ([]Card, error) {
	return []Card{
		{
			Name:     getNameFromPath(path),
			Content:  string(source),
			Filename: getFilename(path),
			Images:   map[string]image.Image{},
		},
	}, nil
}
