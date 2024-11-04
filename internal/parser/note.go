package parser

import "fmt"

// note represents a note parser.
//
// The whole content of the file is returned as a card.
// The is the file name without the extension.
type note struct{}

// newNote returns a new note parser.
func newNote() *note {
	return &note{}
}

// Convert implements the cardParser interface.
func (n *note) parse(path string, source []byte) (Result, error) {
	name := getNameFromPath(path)
	return Result{Cards: []Card{newNoteCard(name, path, source)}}, nil
}

func newNoteCard(name, path string, source []byte) Card {
	content := fmt.Sprintf("# %s\n\n%s", name, string(source))
	return Card{
		Content:  content,
		Fields:   nameFields(name),
		Path:     path,
		Position: sanitizePosition(name),
	}
}

func nameFields(name string) map[string]string {
	return map[string]string{"name": name}
}
