package test

import (
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type ParserCall struct {
	Parser string
	Path   string
	Text   string
	Result parser.Result
	Err    error
}

type ParserCard struct {
	Name       string // for tests
	Content    string
	TemplateID string
	Path       string
	Filename   string
	Position   string
	Images     []parser.Image
	Fields     map[string]mochi.Field
	Equals     bool
}

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
