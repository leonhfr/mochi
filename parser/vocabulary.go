package parser

const (
	vocabularyName          = "vocabulary"
	vocabularyFieldWord     = "word"
	vocabularyFieldExamples = "examples"
	vocabularyFieldNotes    = "notes"
)

type Vocabulary struct{}

func NewVocabulary() *Vocabulary {
	return &Vocabulary{}
}

func (v *Vocabulary) String() string {
	return vocabularyName
}

func (v *Vocabulary) Fields() []string {
	return []string{
		vocabularyFieldWord,
		vocabularyFieldExamples,
		vocabularyFieldNotes,
	}
}
