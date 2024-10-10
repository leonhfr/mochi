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
	Cards  []parser.Card
	Err    error
}

type MockParser struct {
	mock.Mock
}

func NewMockParser(calls []ParserCall) *MockParser {
	m := new(MockParser)
	for _, call := range calls {
		m.
			On("Convert", call.Parser, call.Path, call.Text).
			Return(call.Cards, call.Err)
	}
	return m
}

func (m *MockParser) Convert(parserName, path string, r io.Reader) ([]parser.Card, error) {
	bytes, _ := io.ReadAll(r)
	args := m.Called(parserName, path, string(bytes))
	return args.Get(0).([]parser.Card), args.Error(1)
}
