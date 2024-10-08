package card

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/test"
	"github.com/leonhfr/mochi/mochi"
)

func Test_Parse(t *testing.T) {
	readCalls := []readCall{{
		path:  "/testdata/lorem-ipsum.md",
		bytes: []byte("TEST"),
	}}
	parserCalls := []test.ParserCall{{
		Path:   "/testdata/lorem-ipsum.md",
		Source: []byte("TEST"),
		Cards:  []parser.Card{{Name: "TEST"}},
	}}
	filePaths := []string{"/testdata/lorem-ipsum.md"}
	want := []parser.Card{{Name: "TEST"}}

	r := newMockReader(readCalls)
	p := test.NewMockParser(parserCalls)
	got, err := Parse(r, p, filePaths)
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
		want        []parser.Card
		err         bool
	}{
		{
			name: "read error",
			readCalls: []readCall{{
				path:  "/testdata/lorem-ipsum.md",
				bytes: []byte("TEST"),
				err:   errors.New("ERROR"),
			}},
			path: "/testdata/lorem-ipsum.md",
			err:  true,
		},
		{
			name: "convert error",
			readCalls: []readCall{{
				path:  "/testdata/lorem-ipsum.md",
				bytes: []byte("TEST"),
			}},
			parserCalls: []test.ParserCall{{
				Path:   "/testdata/lorem-ipsum.md",
				Source: []byte("TEST"),
				Err:    errors.New("ERROR"),
			}},
			path: "/testdata/lorem-ipsum.md",
			err:  true,
		},
		{
			name: "success",
			readCalls: []readCall{{
				path:  "/testdata/lorem-ipsum.md",
				bytes: []byte("TEST"),
			}},
			parserCalls: []test.ParserCall{{
				Path:   "/testdata/lorem-ipsum.md",
				Source: []byte("TEST"),
				Cards:  []parser.Card{{Name: "TEST"}},
			}},
			path: "/testdata/lorem-ipsum.md",
			want: []parser.Card{{Name: "TEST"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader(tt.readCalls)
			p := test.NewMockParser(tt.parserCalls)
			got, err := parseFile(r, p, tt.path)
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

func Test_upsertSyncRequests(t *testing.T) {
	filename := "lorem-ipsum.md"
	deckID := "DECK_ID"
	mochiCards := []mochi.Card{
		{
			ID:      "CARD_ID_1",
			Name:    "CARD_TO_UPDATE",
			Content: "OLD_CONTENT",
		},
		{
			ID:      "CARD_ID_2",
			Name:    "CARD_TO_ARCHIVE",
			Content: "CONTENT",
		},
		{
			ID:      "CARD_ID_3",
			Name:    "CARD_TO_KEEP",
			Content: "CONTENT",
		},
	}
	parserCards := []parser.Card{
		{
			Name:     "CARD_TO_UPDATE",
			Content:  "NEW_CONTENT",
			Filename: filename,
		},
		{
			Name:     "CARD_TO_CREATE",
			Content:  "CONTENT",
			Filename: filename,
		},
		{
			Name:     "CARD_TO_KEEP",
			Content:  "CONTENT",
			Filename: filename,
		},
	}

	createWant := []*createCardRequest{
		{
			filename: filename,
			req: mochi.CreateCardRequest{
				Content: "CONTENT",
				DeckID:  "DECK_ID",
				Fields: map[string]mochi.Field{
					"name": {ID: "name", Value: "CARD_TO_CREATE"},
				},
			},
		},
	}
	updateWant := []*updateCardRequest{
		{
			cardID: "CARD_ID_1",
			req:    mochi.UpdateCardRequest{Content: "NEW_CONTENT"},
		},
	}
	archiveWant := []*archiveCardRequest{
		{
			cardID: "CARD_ID_2",
			req:    mochi.UpdateCardRequest{Archived: true},
		},
	}

	gotC, gotU, gotA := upsertSyncRequests(filename, deckID, mochiCards, parserCards)
	assert.Equal(t, createWant, gotC)
	assert.Equal(t, updateWant, gotU)
	assert.Equal(t, archiveWant, gotA)
}

type readCall struct {
	path  string
	bytes []byte
	err   error
}

type mockReader struct {
	mock.Mock
}

func newMockReader(calls []readCall) *mockReader {
	m := new(mockReader)
	for _, call := range calls {
		m.
			On("ReadBytes", call.path).
			Return(call.bytes, call.err)
	}
	return m
}

func (m *mockReader) ReadBytes(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}
