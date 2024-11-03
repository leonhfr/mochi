package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
)

type ParserCall struct {
	Parser string
	Path   string
	Result parser.Result
	Err    error
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
