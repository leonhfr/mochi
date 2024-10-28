package deck

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/test"
)

func Test_Parse(t *testing.T) {
	readCalls := []readCall{{
		path: "/testdata/lorem-ipsum.md",
		text: "TEST",
	}}
	parserCalls := []test.ParserCall{{
		Parser: "note",
		Path:   "/testdata/lorem-ipsum.md",
		Text:   "TEST",
		Result: parser.Result{Cards: []parser.Card{test.NewCard(test.ParserCard{Name: "TEST"})}},
	}}
	filePaths := []string{"/lorem-ipsum.md"}
	want := []Card{
		{card: test.NewCard(test.ParserCard{Name: "TEST"})},
	}

	r := newMockReader(readCalls)
	p := test.NewMockParser(parserCalls)
	got, err := Parse(r, p, "/testdata", "note", filePaths)
	assert.Equal(t, want, got)
	assert.NoError(t, err)
	r.AssertExpectations(t)
	p.AssertExpectations(t)
}

func Test_parseFile(t *testing.T) {
	tests := []struct {
		name        string
		readCalls   []readCall
		parserCalls []test.ParserCall
		path        string
		want        parser.Result
		err         bool
	}{
		{
			name: "read error",
			readCalls: []readCall{{
				path: "/testdata/lorem-ipsum.md",
				text: "TEST",
				err:  errors.New("ERROR"),
			}},
			path: "/lorem-ipsum.md",
			err:  true,
		},
		{
			name: "convert error",
			readCalls: []readCall{{
				path: "/testdata/lorem-ipsum.md",
				text: "TEST",
			}},
			parserCalls: []test.ParserCall{{
				Parser: "note",
				Path:   "/testdata/lorem-ipsum.md",
				Text:   "TEST",
				Err:    errors.New("ERROR"),
			}},
			path: "/lorem-ipsum.md",
			err:  true,
		},
		{
			name: "success",
			readCalls: []readCall{{
				path: "/testdata/lorem-ipsum.md",
				text: "TEST",
			}},
			parserCalls: []test.ParserCall{{
				Parser: "note",
				Path:   "/testdata/lorem-ipsum.md",
				Text:   "TEST",
				Result: parser.Result{Cards: []parser.Card{test.NewCard(test.ParserCard{Name: "TEST"})}},
			}},
			path: "/lorem-ipsum.md",
			want: parser.Result{Cards: []parser.Card{test.NewCard(test.ParserCard{Name: "TEST"})}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader(tt.readCalls)
			p := test.NewMockParser(tt.parserCalls)
			got, err := parseFile(r, p, "/testdata", "note", tt.path)
			assert.Equal(t, tt.want, got)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			r.AssertExpectations(t)
			p.AssertExpectations(t)
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
