package parser

var (
	_ Parser = &Headings{}
	_ Parser = &Note{}
	_ Parser = &Vocabulary{}
)
