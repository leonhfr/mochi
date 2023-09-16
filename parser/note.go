package parser

const (
	noteName         = "note"
	noteFieldTitle   = "title"
	noteFieldContent = "content"
)

type Note struct{}

func NewNote() *Note {
	return &Note{}
}

func (n *Note) String() string {
	return noteName
}

func (n *Note) Fields() []string {
	return []string{noteFieldTitle, noteFieldContent}
}
