package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/mochi"
)

type ParserCall struct {
	Parser string
	Path   string
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
			On("Parse", mock.Anything, call.Parser, call.Path).
			Return(call.Result, call.Err)
	}
	return m
}

func (m *MockParser) Parse(reader parser.Reader, parserName, path string) (parser.Result, error) {
	args := m.Called(reader, parserName, path)
	return args.Get(0).(parser.Result), args.Error(1)
}
