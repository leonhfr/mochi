package deck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/test"
	"github.com/leonhfr/mochi/mochi"
)

func Test_Sync(t *testing.T) {
	tests := []struct {
		name     string
		client   test.Mochi
		config   test.Config
		lockfile test.Lockfile
		path     string
		want     string
		err      error
	}{
		{
			name: "should create the expected decks",
			client: test.Mochi{
				CreateDeck: []test.MochiCreateDeck{
					{
						Req:  mochi.CreateDeckRequest{Name: "DECK_DATA_NAME", ParentID: "DECK_TEST_ID"},
						Deck: mochi.Deck{ID: "DECK_DATA_ID"},
					},
				},
			},
			config: test.Config{
				Deck: []test.ConfigDeck{
					{Path: "/test/data", Deck: config.Deck{Name: "DECK_DATA_NAME"}, OK: true},
				},
			},
			lockfile: test.Lockfile{
				Lock: 1,
				GetDeck: []test.LockfileDeck{
					{Path: "/test/data"},
					{Path: "/test", DeckID: "DECK_TEST_ID", OK: true},
				},
				SetDeck: []test.LockfileSetDeck{
					{ID: "DECK_DATA_ID", ParentID: "DECK_TEST_ID", Path: "/test/data", Name: "DECK_DATA_NAME"},
				},
			},
			path: "/test/data",
			want: "DECK_DATA_ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := test.NewMockMochi(tt.client)
			config := test.NewMockConfig(tt.config)
			lf := test.NewMockLockfile(tt.lockfile)
			got, err := Sync(context.Background(), client, config, lf, tt.path)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
			client.AssertExpectations(t)
			config.AssertExpectations(t)
			lf.AssertExpectations(t)
		})
	}
}

func Test_createDeck(t *testing.T) {
	tests := []struct {
		name     string
		client   []test.MochiCreateDeck
		lockfile []test.LockfileSetDeck
		parentID string
		path     string
		deckName string
		want     string
		err      error
	}{
		{
			name: "api error",
			client: []test.MochiCreateDeck{
				{Req: mochi.CreateDeckRequest{Name: "DECK_NAME", ParentID: "PARENT_ID"}, Err: test.ErrMochi},
			},
			parentID: "PARENT_ID",
			deckName: "DECK_NAME",
			err:      test.ErrMochi,
		},
		{
			name: "success",
			client: []test.MochiCreateDeck{
				{Req: mochi.CreateDeckRequest{Name: "DECK_NAME", ParentID: "PARENT_ID"}, Deck: mochi.Deck{ID: "DECK_ID"}},
			},
			lockfile: []test.LockfileSetDeck{
				{ID: "DECK_ID", ParentID: "PARENT_ID", Path: "/test/data", Name: "DECK_NAME"},
			},
			parentID: "PARENT_ID",
			path:     "/test/data",
			deckName: "DECK_NAME",
			want:     "DECK_ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := test.NewMockMochi(test.Mochi{CreateDeck: tt.client})
			lf := test.NewMockLockfile(test.Lockfile{SetDeck: tt.lockfile})
			got, err := createDeck(context.Background(), client, lf, tt.parentID, tt.path, tt.deckName)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
			client.AssertExpectations(t)
			lf.AssertExpectations(t)
		})
	}
}

func Test_updateDeckName(t *testing.T) {
	tests := []struct {
		name     string
		client   []test.MochiUpdateDeck
		lockfile []test.LockfileUpdateDeckName
		deckID   string
		deckName string
		err      error
	}{
		{
			name: "api error",
			client: []test.MochiUpdateDeck{
				{ID: "DECK_ID", Req: mochi.UpdateDeckRequest{Name: "DECK_NAME"}, Err: test.ErrMochi},
			},
			deckID:   "DECK_ID",
			deckName: "DECK_NAME",
			err:      test.ErrMochi,
		},
		{
			name: "success",
			client: []test.MochiUpdateDeck{
				{ID: "DECK_ID", Req: mochi.UpdateDeckRequest{Name: "DECK_NAME"}},
			},
			lockfile: []test.LockfileUpdateDeckName{
				{ID: "DECK_ID", Name: "DECK_NAME"},
			},
			deckID:   "DECK_ID",
			deckName: "DECK_NAME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := test.NewMockMochi(test.Mochi{UpdateDeck: tt.client})
			lf := test.NewMockLockfile(test.Lockfile{UpdateDeckName: tt.lockfile})
			err := updateDeckName(context.Background(), client, lf, tt.deckID, tt.deckName)
			assert.Equal(t, tt.err, err)
			client.AssertExpectations(t)
			lf.AssertExpectations(t)
		})
	}
}

func Test_getDeckName(t *testing.T) {
	tests := []struct {
		name  string
		calls []test.ConfigDeck
		path  string
		want  string
	}{
		{
			name: "config deck has name",
			calls: []test.ConfigDeck{
				{Path: "/test/data", Deck: config.Deck{Name: "DECK_NAME"}, OK: true},
			},
			path: "/test/data",
			want: "DECK_NAME",
		},
		{
			name: "config deck has empty name",
			calls: []test.ConfigDeck{
				{Path: "/test/data", Deck: config.Deck{}, OK: true},
			},
			path: "/test/data",
			want: "data",
		},
		{
			name: "no config deck",
			calls: []test.ConfigDeck{
				{Path: "/test/data", OK: false},
			},
			path: "/test/data",
			want: "data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := test.NewMockConfig(test.Config{Deck: tt.calls})
			got := getDeckName(cfg, tt.path)
			assert.Equal(t, tt.want, got)
			cfg.AssertExpectations(t)
		})
	}
}

func Test_getStack(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		calls  []test.LockfileDeck
		deckID string
		stack  []string
	}{
		{
			name:  "path is root",
			path:  "/",
			stack: []string{"/"},
		},
		{
			name: "recursive to top level directory",
			path: "/test/data/notes",
			calls: []test.LockfileDeck{
				{Path: "/test/data"},
				{Path: "/test"},
			},
			stack: []string{"/test/data/notes", "/test/data", "/test"},
		},
		{
			name: "recursive to existing",
			path: "/test/data/notes",
			calls: []test.LockfileDeck{
				{Path: "/test/data"},
				{Path: "/test", DeckID: "DECK_ID", OK: true},
			},
			deckID: "DECK_ID",
			stack:  []string{"/test/data/notes", "/test/data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lf := test.NewMockLockfile(test.Lockfile{GetDeck: tt.calls})
			deckID, stack := getStack(lf, tt.path)
			assert.Equal(t, tt.deckID, deckID)
			assert.Equal(t, tt.stack, stack)
			lf.AssertExpectations(t)
		})
	}
}

func Test_isTopLevelDirectory(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "/", want: true},
		{path: "/testdata", want: true},
		{path: "/testdata/lorem-ipsum"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isTopLevelDirectory(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}
