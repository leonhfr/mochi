package card

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/converter"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/test"
)

func Test_Parse(t *testing.T) {
	parserCalls := []test.ParserCall{{
		Parser: "note",
		Path:   "/testdata/lorem-ipsum.md",
		Result: parser.Result{Cards: []parser.Card{{Content: "TEST"}}},
	}}
	converterCalls := []test.ConverterCall{{
		Path:   "/testdata/lorem-ipsum.md",
		Source: "TEST",
		Result: converter.Result{Markdown: "TEST"},
	}}
	filePaths := []string{"/lorem-ipsum.md"}
	want := []Card{
		{Card: parser.Card{Content: "TEST"}},
	}

	p := test.NewMockParser(parserCalls)
	c := test.NewMockConverter(converterCalls)
	got, err := Parse(nil, p, c, "/testdata", "note", filePaths)
	assert.Equal(t, want, got)
	assert.NoError(t, err)
	p.AssertExpectations(t)
	c.AssertExpectations(t)
}

func Test_parseFile(t *testing.T) {
	tests := []struct {
		name        string
		parserCalls []test.ParserCall
		path        string
		wantDeck    string
		wantCards   []parser.Card
		err         bool
	}{
		{
			name: "convert error",
			parserCalls: []test.ParserCall{{
				Parser: "note",
				Path:   "/testdata/lorem-ipsum.md",
				Err:    errors.New("ERROR"),
			}},
			path: "/testdata/lorem-ipsum.md",
			err:  true,
		},
		{
			name: "success",
			parserCalls: []test.ParserCall{{
				Parser: "note",
				Path:   "/testdata/lorem-ipsum.md",
				Result: parser.Result{Cards: []parser.Card{{Content: "TEST"}}},
			}},
			path:      "/testdata/lorem-ipsum.md",
			wantCards: []parser.Card{{Content: "TEST"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := test.NewMockParser(tt.parserCalls)
			gotDeck, gotCards, err := parseFile(nil, p, "note", tt.path)
			assert.Equal(t, tt.wantDeck, gotDeck)
			assert.Equal(t, tt.wantCards, gotCards)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			p.AssertExpectations(t)
		})
	}
}
