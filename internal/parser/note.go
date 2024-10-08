package parser

// Note represents a note parser.
//
// The whole content of the file is returned as a card.
// The is the file name without the extension.
type Note struct{}

// NewNote returns a new note parser.
func NewNote() *Note {
	return &Note{}
}

// Convert implements the Parser interface.
func (n *Note) Convert(path string, source []byte) ([]Card, error) {
	return []Card{
		{
			Name:    getNameFromPath(path),
			Content: string(source),
		},
	}, nil
}
