package lock

import (
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name        string
		target      string
		path        string
		fileContent string
		fileError   error
		wantPath    string
		wantData    map[string]Deck
		err         bool
	}{
		{
			name:      "no lockfile found",
			target:    "testdata",
			path:      "testdata/mochi-lock.json",
			fileError: fs.ErrNotExist,
			wantPath:  "testdata/mochi-lock.json",
			wantData:  map[string]Deck{},
		},
		{
			name:        "lockfile found",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			fileContent: `{"DECK_ID":{"path":"DECK_PATH","name":"DECK_NAME"}}`,
			wantPath:    "testdata/mochi-lock.json",
			wantData:    map[string]Deck{"DECK_ID": {Path: "DECK_PATH", Name: "DECK_NAME"}},
		},
		{
			name:        "missing path and name",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			fileContent: `{"DECK_ID":{}}`,
			err:         true,
		},
		{
			name:        "missing card filename",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			fileContent: `{"DECK_ID":{"path":"DECK_PATH","name":"DECK_NAME","cards":{"CARD_ID":{}}}}`,
			err:         true,
		},
		{
			name:        "both path and virtual",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			fileContent: `{"DECK_ID":{"path":"DECK_PATH","name":"DECK_NAME","virtual":true,"cards":{}}}`,
			err:         true,
		},
		{
			name:        "neither path or virtual",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			fileContent: `{"DECK_ID":{"name":"DECK_NAME","cards":{}}}`,
			err:         true,
		},
		{
			name:        "virtual deck",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			fileContent: `{"DECK_ID":{"name":"DECK_NAME","virtual":true,"cards":{}}}`,
			wantPath:    "testdata/mochi-lock.json",
			wantData:    map[string]Deck{"DECK_ID": {Name: "DECK_NAME", Virtual: true, Cards: map[string]Card{}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := new(mockFile)
			rw.On("Read", tt.path).Return(tt.fileContent, tt.fileError)

			got, err := Parse(rw, tt.target)

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPath, got.path)
				assert.Equal(t, tt.wantData, got.decks)
			}
			rw.AssertExpectations(t)
		})
	}
}

func Test_Lock_DeckFromPath(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]Deck
		path   string
		deckID string
		deck   Deck
		ok     bool
	}{
		{
			name:   "deck found",
			data:   map[string]Deck{"DECK_ID": {Name: "DECK_NAME", Path: "/lorem-ipsum"}},
			path:   "/lorem-ipsum",
			deckID: "DECK_ID",
			deck:   Deck{Name: "DECK_NAME", Path: "/lorem-ipsum"},
			ok:     true,
		},
		{
			name: "deck not found",
			data: map[string]Deck{"DECK_ID": {Name: "DECK_NAME", Path: "/lorem-ipsum"}},
			path: "/sed-interdum-libero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{decks: tt.data}
			deckID, deck, ok := lock.DeckFromPath(tt.path)
			assert.Equal(t, tt.deckID, deckID)
			assert.Equal(t, tt.deck, deck)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Lock_SetDeck(t *testing.T) {
	deckID, parentID, path, name := "DECK_ID", "PARENT_DECK_ID", "/deck", "Deck"
	want := map[string]Deck{
		deckID: {ParentID: parentID, Path: path, Name: name, Cards: map[string]Card{}},
	}
	lock := &Lock{decks: make(map[string]Deck)}
	lock.SetDeck(deckID, parentID, path, name)
	assert.Equal(t, lock.decks, want)
	assert.True(t, lock.updated)
}

func Test_Lock_UpdateDeck(t *testing.T) {
	deckID, parentID, path, name := "DECK_ID", "PARENT_DECK_ID", "/deck", "Deck"
	want := "Updated deck name"
	lock := &Lock{decks: map[string]Deck{deckID: {ParentID: parentID, Path: path, Name: name}}}
	lock.UpdateDeck(deckID, want)
	assert.Equal(t, lock.decks[deckID].Name, want)
	assert.True(t, lock.updated)
}

func Test_Lock_Card(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]Deck
		deckID string
		cardID string
		want   Card
		ok     bool
	}{
		{
			name:   "deck does not exist",
			deckID: "DECK_ID",
			cardID: "CARD_ID",
		},
		{
			name:   "card does not exist",
			data:   map[string]Deck{"DECK_ID": {}},
			deckID: "DECK_ID",
			cardID: "CARD_ID",
		},
		{
			name: "card found",
			data: map[string]Deck{"DECK_ID": {
				Cards: map[string]Card{"CARD_ID": {Filename: "lorem-ipsum.md"}},
			}},
			deckID: "DECK_ID",
			cardID: "CARD_ID",
			want:   Card{Filename: "lorem-ipsum.md"},
			ok:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{decks: tt.data}
			got, ok := lock.Card(tt.deckID, tt.cardID)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Lock_SetCard(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]Deck
		deckID   string
		cardID   string
		filename string
		want     map[string]Deck
		err      bool
	}{
		{
			name:     "deck does not exist",
			data:     map[string]Deck{},
			deckID:   "DECK_ID",
			cardID:   "CARD_ID",
			filename: "/lorem-ipsum.md",
			want:     map[string]Deck{},
			err:      true,
		},
		{
			name: "rewrite when card already exists",
			data: map[string]Deck{"DECK_ID": {
				Cards: map[string]Card{"CARD_ID": {Filename: "/old.md"}},
			}},
			deckID:   "DECK_ID",
			cardID:   "CARD_ID",
			filename: "/lorem-ipsum.md",
			want: map[string]Deck{"DECK_ID": {
				Cards: map[string]Card{"CARD_ID": {Filename: "/lorem-ipsum.md"}},
			}},
		},
		{
			name: "card set",
			data: map[string]Deck{"DECK_ID": {
				Cards: map[string]Card{},
			}},
			deckID:   "DECK_ID",
			cardID:   "CARD_ID",
			filename: "/lorem-ipsum.md",
			want: map[string]Deck{"DECK_ID": {
				Cards: map[string]Card{"CARD_ID": {
					Filename: "/lorem-ipsum.md",
				}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{decks: tt.data}
			err := lock.SetCard(tt.deckID, tt.cardID, tt.filename)
			assert.Equal(t, tt.want, lock.decks)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockFile struct {
	mock.Mock
}

func (m *mockFile) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}

func (m *mockFile) Write(p string) (io.WriteCloser, error) {
	args := m.Mock.Called(p)
	wc := writeCloser{&strings.Builder{}}
	return wc, args.Error(1)
}

type writeCloser struct {
	*strings.Builder
}

func (writeCloser) Close() error { return nil }
