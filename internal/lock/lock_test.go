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
			name:        "bad config",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			exists:      true,
			fileContent: `{"DECK_ID":{}}`,
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
