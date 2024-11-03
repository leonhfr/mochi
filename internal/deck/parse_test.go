package deck

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/test"
)

func Test_Parse(t *testing.T) {
	parserCalls := []test.ParserCall{{
		Parser: "note",
		Path:   "/testdata/lorem-ipsum.md",
		Result: parser.Result{Cards: []parser.Card{{Content: "TEST"}}},
	}}
	filePaths := []string{"/lorem-ipsum.md"}
	want := []Card{
		{card: parser.Card{Content: "TEST"}},
	}

	p := test.NewMockParser(parserCalls)
	got, err := Parse(nil, p, "/testdata", "note", filePaths)
	assert.Equal(t, want, got)
	assert.NoError(t, err)
	p.AssertExpectations(t)
}

func Test_parseFile(t *testing.T) {
	tests := []struct {
		name        string
		parserCalls []test.ParserCall
		path        string
		want        parser.Result
		err         bool
	}{
		{
			name: "convert error",
			parserCalls: []test.ParserCall{{
				Parser: "note",
				Path:   "/testdata/lorem-ipsum.md",
				Err:    errors.New("ERROR"),
			}},
			path: "/lorem-ipsum.md",
			err:  true,
		},
		{
			name: "success",
			parserCalls: []test.ParserCall{{
				Parser: "note",
				Path:   "/testdata/lorem-ipsum.md",
				Result: parser.Result{Cards: []parser.Card{{Content: "TEST"}}},
			}},
			path: "/lorem-ipsum.md",
			want: parser.Result{Cards: []parser.Card{{Content: "TEST"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := test.NewMockParser(tt.parserCalls)
			got, err := parseFile(nil, p, "/testdata", "note", tt.path)
			assert.Equal(t, tt.want, got)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			p.AssertExpectations(t)
		})
	}
}
