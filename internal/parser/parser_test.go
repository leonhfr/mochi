package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Parser_Convert(t *testing.T) {
	mockPath := "/testdata/lorem-ipsum/Lorem ipsum.md"
	mockSource := "# Title 1\nParagraph.\n"
	mockCards := []Card{{
		Name:     "Lorem ipsum",
		Content:  "# Title 1\nParagraph.\n",
		Filename: "Lorem ipsum.md",
	}}

	tests := []struct {
		name    string
		parser  string
		path    string
		parser0 []cardParserCall
		parser1 []cardParserCall
		parser2 []cardParserCall
		source  string
		want    []Card
	}{
		{
			name:   "frontmatter skip",
			path:   mockPath,
			source: "---\nmochi-skip: true\n---\n" + mockSource,
		},
		{
			name:    "frontmatter overwrites parser",
			parser2: []cardParserCall{{path: mockPath, source: mockSource, cards: mockCards}},
			parser:  "parser1",
			path:    mockPath,
			source:  "---\nmochi-parser: parser2\n---\n" + mockSource,
			want:    mockCards,
		},
		{
			name:    "specific parser",
			parser1: []cardParserCall{{path: mockPath, source: mockSource, cards: mockCards}},
			parser:  "parser1",
			path:    mockPath,
			source:  mockSource,
			want:    mockCards,
		},
		{
			name:    "default parser",
			parser0: []cardParserCall{{path: mockPath, source: mockSource, cards: mockCards}},
			path:    mockPath,
			source:  mockSource,
			want:    mockCards,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			got, err := parser.Convert(tt.parser, tt.path, strings.NewReader(tt.source))
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
			p0.AssertExpectations(t)
			p1.AssertExpectations(t)
			p2.AssertExpectations(t)
		})
	}
}

type cardParserCall struct {
	path   string
	source string
	cards  []Card
	err    error
}

type mockCardParser struct {
	mock.Mock
}

func newMockCardParser(calls []cardParserCall) *mockCardParser {
	m := new(mockCardParser)
	for _, call := range calls {
		m.
			On("convert", call.path, []byte(call.source)).
			Return(call.cards, call.err)
	}
	return m
}

func (m *mockCardParser) convert(path string, source []byte) ([]Card, error) {
	args := m.Called(path, source)
	return args.Get(0).([]Card), args.Error(1)
}
