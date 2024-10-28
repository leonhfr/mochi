package test

import (
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
)

type ParserCall struct {
	Parser string
	Path   string
	Text   string
	Result parser.Result
	Err    error
}

type ParserCard struct {
	Name     string
	Content  string
	Path     string
	Filename string
	Position string
	Images   []parser.Image
}

type parserCard struct {
	name     string
	content  string
	path     string
	filename string
	position string
	images   []parser.Image
}

func NewCard(card ParserCard) parser.Card {
	return &parserCard{
		name:     card.Name,
		content:  card.Content,
		path:     card.Path,
		filename: card.Filename,
		position: card.Position,
		images:   card.Images,
	}
}

var _ parser.Card = (*parserCard)(nil)

func (c *parserCard) Name() string           { return c.name }
func (c *parserCard) Content() string        { return c.content }
func (c *parserCard) Path() string           { return c.path }
func (c *parserCard) Filename() string       { return c.filename }
func (c *parserCard) Position() string       { return c.position }
func (c *parserCard) Images() []parser.Image { return c.images }

type MockParser struct {
	mock.Mock
}

func NewMockParser(calls []ParserCall) *MockParser {
	m := new(MockParser)
	for _, call := range calls {
		m.
			On("Convert", call.Parser, call.Path, call.Text).
			Return(call.Result, call.Err)
	}
	return m
}

func (m *MockParser) Convert(parserName, path string, r io.Reader) (parser.Result, error) {
	bytes, _ := io.ReadAll(r)
	args := m.Called(parserName, path, string(bytes))
	return args.Get(0).(parser.Result), args.Error(1)
}
