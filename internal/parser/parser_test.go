package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Parser_Convert(t *testing.T) {
	mockPath := "/testdata/lorem-ipsum/Lorem ipsum.md"
	mockSource := "# Title 1\nParagraph.\n"
	mockCards := []Card{
		{
			Content: "# Title 1\nParagraph.\n",
		},
	}

	tests := []struct {
		name    string
		parser  string
		path    string
		parser0 []cardParserCall
		parser1 []cardParserCall
		parser2 []cardParserCall
		source  string
		want    Result
	}{
		{
			name:   "frontmatter skip",
			path:   mockPath,
			source: "---\nmochi-skip: true\n---\n" + mockSource,
		},
		{
			name:    "frontmatter overwrites parser",
			parser2: []cardParserCall{{path: mockPath, source: mockSource, result: Result{Cards: mockCards}}},
			parser:  "parser1",
			path:    mockPath,
			source:  "---\nmochi-parser: parser2\n---\n" + mockSource,
			want:    Result{Cards: mockCards},
		},
		{
			name:    "specific parser",
			parser1: []cardParserCall{{path: mockPath, source: mockSource, result: Result{Cards: mockCards}}},
			parser:  "parser1",
			path:    mockPath,
			source:  mockSource,
			want:    Result{Cards: mockCards},
		},
		{
			name:    "default parser",
			parser0: []cardParserCall{{path: mockPath, source: mockSource, result: Result{Cards: mockCards}}},
			path:    mockPath,
			source:  mockSource,
			want:    Result{Cards: mockCards},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader([]readCall{{path: tt.path, text: tt.source}})
			p0 := newMockCardParser(tt.parser0)
			p1 := newMockCardParser(tt.parser1)
			p2 := newMockCardParser(tt.parser2)
			parser := &Parser{
				cardParser: p0,
				parsers: map[string]cardParser{
					"parser1": p1,
					"parser2": p2,
				},
			}

			got, err := parser.Parse(r, tt.parser, tt.path)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
			p0.AssertExpectations(t)
			p1.AssertExpectations(t)
			p2.AssertExpectations(t)
		})
	}
}

type readCall struct {
	path string
	text string
	err  error
}

type mockReader struct {
	mock.Mock
}

func newMockReader(calls []readCall) *mockReader {
	m := new(mockReader)
	for _, call := range calls {
		m.
			On("Read", call.path).
			Return(call.text, call.err)
	}
	return m
}

func (m *mockReader) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}

type cardParserCall struct {
	path   string
	source string
	result Result
	err    error
}

type mockCardParser struct {
	mock.Mock
}

func newMockCardParser(calls []cardParserCall) *mockCardParser {
	m := new(mockCardParser)
	for _, call := range calls {
		m.
			On("parse", call.path, []byte(call.source)).
			Return(call.result, call.err)
	}
	return m
}

func (m *mockCardParser) parse(path string, source []byte) (Result, error) {
	args := m.Called(path, source)
	return args.Get(0).(Result), args.Error(1)
}
