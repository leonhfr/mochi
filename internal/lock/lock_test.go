package lock

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/mochi"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name        string
		target      string
		path        string
		exists      bool
		fileContent string
		fileError   error
		wantPath    string
		wantData    lockData
		err         bool
	}{
		{
			name:     "no lockfile found",
			target:   "testdata",
			path:     "testdata/mochi-lock.json",
			wantPath: "testdata/mochi-lock.json",
			wantData: lockData{},
		},
		{
			name:        "lockfile found",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			exists:      true,
			fileContent: `{"DECK_ID":{"path":"DECK_PATH","name":"DECK_NAME"}}`,
			wantPath:    "testdata/mochi-lock.json",
			wantData:    lockData{"DECK_ID": Deck{Path: "DECK_PATH", Name: "DECK_NAME"}},
		},
		{
			name:        "missing path and name",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			exists:      true,
			fileContent: `{"DECK_ID":{}}`,
			err:         true,
		},
		{
			name:        "missing card filename",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			exists:      true,
			fileContent: `{"DECK_ID":{"path":"DECK_PATH","name":"DECK_NAME","cards":{"CARD_ID":{}}}}`,
			err:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := new(mockFile)
			rw.On("Exists", tt.path).Return(tt.exists)
			if tt.exists {
				rw.On("Read", tt.path).Return(tt.fileContent, tt.fileError)
			}

			got, err := Parse(rw, tt.target)

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.wantPath, got.path)
				assert.Equal(t, tt.wantData, got.data)
				assert.NoError(t, err)
			}
			rw.AssertExpectations(t)
		})
	}
}

func Test_Lock_CleanDecks(t *testing.T) {
	tests := []struct {
		name    string
		decks   []mochi.Deck
		data    lockData
		want    lockData
		updated bool
	}{
		{
			name: "should not modify the lock",
			decks: []mochi.Deck{
				{ID: "DECK_ID_2", Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			data: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
				"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
			},
			want: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
				"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
			},
			updated: false,
		},
		{
			name: "should remove decks that are not in the slice",
			decks: []mochi.Deck{
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			data: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
				"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
			},
			want: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
			},
			updated: true,
		},
		{
			name: "should remove decks whose parent id have changed",
			decks: []mochi.Deck{
				{ID: "DECK_ID_2", Name: "DECK_NAME_2", ParentID: ""},
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			data: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
				"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
			},
			want: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
			},
			updated: true,
		},
		{
			name: "should update the deck name",
			decks: []mochi.Deck{
				{ID: "DECK_ID_2", Name: "NEW_DECK_NAME_2", ParentID: "DECK_ID_1"},
				{ID: "DECK_ID_1", Name: "DECK_NAME_1", ParentID: ""},
			},
			data: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
				"DECK_ID_2": {Name: "DECK_NAME_2", ParentID: "DECK_ID_1"},
			},
			want: lockData{
				"DECK_ID_1": {Name: "DECK_NAME_1", ParentID: ""},
				"DECK_ID_2": {Name: "NEW_DECK_NAME_2", ParentID: "DECK_ID_1"},
			},
			updated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{data: tt.data}
			lock.CleanDecks(tt.decks)
			assert.Equal(t, lock.data, tt.want)
			assert.Equal(t, lock.updated, tt.updated)
		})
	}
}

func Test_Lock_CleanCards(t *testing.T) {
	tests := []struct {
		name    string
		deckID  string
		cardIDs []string
		data    lockData
		want    lockData
		updated bool
	}{
		{
			name:    "should not modify the lock if deck is not found",
			deckID:  "DECK_ID_1",
			cardIDs: []string{"CARD_ID_1", "CARD_ID_2", "CARD_ID_3"},
			data: lockData{
				"DECK_ID_2": {Cards: map[string]Card{
					"CARD_ID_1": {},
					"CARD_ID_2": {},
					"CARD_ID_3": {},
				}},
			},
			want: lockData{
				"DECK_ID_2": {Cards: map[string]Card{
					"CARD_ID_1": {},
					"CARD_ID_2": {},
					"CARD_ID_3": {},
				}},
			},
			updated: false,
		},
		{
			name:    "should not modify the lock if cards are present",
			deckID:  "DECK_ID",
			cardIDs: []string{"CARD_ID_1", "CARD_ID_2", "CARD_ID_3"},
			data: lockData{
				"DECK_ID": {Cards: map[string]Card{
					"CARD_ID_1": {},
					"CARD_ID_2": {},
					"CARD_ID_3": {},
				}},
			},
			want: lockData{
				"DECK_ID": {Cards: map[string]Card{
					"CARD_ID_1": {},
					"CARD_ID_2": {},
					"CARD_ID_3": {},
				}},
			},
			updated: false,
		},
		{
			name:    "should remove cards that are not in the slice",
			deckID:  "DECK_ID",
			cardIDs: []string{"CARD_ID_1", "CARD_ID_2"},
			data: lockData{
				"DECK_ID": {Cards: map[string]Card{
					"CARD_ID_1": {},
					"CARD_ID_2": {},
					"CARD_ID_3": {},
				}},
			},
			want: lockData{
				"DECK_ID": {Cards: map[string]Card{
					"CARD_ID_1": {},
					"CARD_ID_2": {},
				}},
			},
			updated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{data: tt.data}
			lock.CleanCards(tt.deckID, tt.cardIDs)
			assert.Equal(t, lock.data, tt.want)
			assert.Equal(t, lock.updated, tt.updated)
		})
	}
}

func Test_Lock_GetDeck(t *testing.T) {
	tests := []struct {
		name   string
		data   lockData
		path   string
		deckID string
		deck   Deck
		ok     bool
	}{
		{
			name:   "deck found",
			data:   lockData{"DECK_ID": {Name: "DECK_NAME", Path: "/lorem-ipsum"}},
			path:   "/lorem-ipsum",
			deckID: "DECK_ID",
			deck:   Deck{Name: "DECK_NAME", Path: "/lorem-ipsum"},
			ok:     true,
		},
		{
			name: "deck not found",
			data: lockData{"DECK_ID": {Name: "DECK_NAME", Path: "/lorem-ipsum"}},
			path: "/sed-interdum-libero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{data: tt.data}
			deckID, deck, ok := lock.GetDeck(tt.path)
			assert.Equal(t, tt.deckID, deckID)
			assert.Equal(t, tt.deck, deck)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Lock_SetDeck(t *testing.T) {
	deckID, parentID, path, name := "DECK_ID", "PARENT_DECK_ID", "/deck", "Deck"
	want := lockData{
		deckID: Deck{ParentID: parentID, Path: path, Name: name, Cards: map[string]Card{}},
	}
	lock := &Lock{data: make(lockData)}
	lock.SetDeck(deckID, parentID, path, name)
	assert.Equal(t, lock.data, want)
	assert.True(t, lock.updated)
}

func Test_Lock_UpdateDeckName(t *testing.T) {
	deckID, parentID, path, name := "DECK_ID", "PARENT_DECK_ID", "/deck", "Deck"
	want := "Updated deck name"
	lock := &Lock{data: lockData{deckID: Deck{ParentID: parentID, Path: path, Name: name}}}
	lock.UpdateDeckName(deckID, want)
	assert.Equal(t, lock.data[deckID].Name, want)
	assert.True(t, lock.updated)
}

func Test_Lock_GetCard(t *testing.T) {
	tests := []struct {
		name   string
		data   lockData
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
			data:   lockData{"DECK_ID": Deck{}},
			deckID: "DECK_ID",
			cardID: "CARD_ID",
		},
		{
			name: "card found",
			data: lockData{"DECK_ID": Deck{
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
			lock := &Lock{data: tt.data}
			got, ok := lock.GetCard(tt.deckID, tt.cardID)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Lock_SetCard(t *testing.T) {
	tests := []struct {
		name     string
		data     lockData
		deckID   string
		cardID   string
		filename string
		want     lockData
		err      bool
	}{
		{
			name:     "deck does not exist",
			data:     lockData{},
			deckID:   "DECK_ID",
			cardID:   "CARD_ID",
			filename: "/lorem-ipsum.md",
			want:     lockData{},
			err:      true,
		},
		{
			name: "card already exists",
			data: lockData{"DECK_ID": Deck{
				Cards: map[string]Card{"CARD_ID": {}},
			}},
			deckID:   "DECK_ID",
			cardID:   "CARD_ID",
			filename: "/lorem-ipsum.md",
			want: lockData{"DECK_ID": Deck{
				Cards: map[string]Card{"CARD_ID": {}},
			}},
		},
		{
			name: "card set",
			data: lockData{"DECK_ID": Deck{
				Cards: map[string]Card{},
			}},
			deckID:   "DECK_ID",
			cardID:   "CARD_ID",
			filename: "/lorem-ipsum.md",
			want: lockData{"DECK_ID": Deck{
				Cards: map[string]Card{"CARD_ID": {Filename: "/lorem-ipsum.md"}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &Lock{data: tt.data}
			err := lock.SetCard(tt.deckID, tt.cardID, tt.filename)
			assert.Equal(t, tt.want, lock.data)
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

func (m *mockFile) Exists(p string) bool {
	args := m.Mock.Called(p)
	return args.Bool(0)
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
