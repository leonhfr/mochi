package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
)

type ParserCall struct {
	Path   string
	Source []byte
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
			On("Convert", call.Path, call.Source).
			Return(call.Cards, call.Err)
	}
	return m
}

func (m *MockParser) Convert(path string, source []byte) ([]parser.Card, error) {
	args := m.Called(path, source)
	return args.Get(0).([]parser.Card), args.Error(1)
}
