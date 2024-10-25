package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/lock"
)

type Lockfile struct {
	Lock         int
	Deck         []LockfileDeck
	Decks        []map[string]lock.Deck
	DeckFromPath []LockfileDeckFromPath
	SetDeck      []LockfileSetDeck
	DeleteDeck   []string
	UpdateDeck   []LockfileUpdateDeck
	DeleteCard   []LockfileDeleteCard
}

type LockfileDeck struct {
	DeckID string
	Deck   lock.Deck
	OK     bool
}

type LockfileDeckFromPath struct {
	Path   string
	DeckID string
	Deck   lock.Deck
	OK     bool
}

type LockfileSetDeck struct {
	ID       string
	ParentID string
	Path     string
	Name     string
}

type LockfileUpdateDeck struct {
	ID   string
	Name string
}

type LockfileDeleteCard struct {
	DeckID string
	CardID string
}

func NewMockLockfile(calls Lockfile) *MockLockfile {
	lf := new(MockLockfile)
	for i := 0; i < calls.Lock; i++ {
		lf.On("Lock").Return()
		lf.On("Unlock").Return()
	}
	for _, call := range calls.Deck {
		lf.On("Deck", call.DeckID).Return(call.Deck, call.OK)
	}
	for _, call := range calls.DeckFromPath {
		lf.On("DeckFromPath", call.Path).Return(call.DeckID, call.Deck, call.OK)
	}
	for _, call := range calls.Decks {
		lf.On("Decks").Return(call)
	}
	for _, call := range calls.SetDeck {
		lf.On("SetDeck", call.ID, call.ParentID, call.Path, call.Name).Return()
	}
	for _, call := range calls.UpdateDeck {
		lf.On("UpdateDeck", call.ID, call.Name).Return()
	}
	for _, call := range calls.DeleteDeck {
		lf.On("DeleteDeck", call).Return()
	}
	for _, call := range calls.DeleteCard {
		lf.On("DeleteCard", call.DeckID, call.CardID).Return()
	}
	return lf
}

type MockLockfile struct {
	mock.Mock
}

func (m *MockLockfile) Lock() {
	m.Called()
}

func (m *MockLockfile) Unlock() {
	m.Called()
}

func (m *MockLockfile) Deck(id string) (lock.Deck, bool) {
	args := m.Called(id)
	return args.Get(0).(lock.Deck), args.Bool(1)
}

func (m *MockLockfile) Decks() map[string]lock.Deck {
	args := m.Called()
	return args.Get(0).(map[string]lock.Deck)
}

func (m *MockLockfile) DeckFromPath(path string) (string, lock.Deck, bool) {
	args := m.Called(path)
	return args.String(0), args.Get(1).(lock.Deck), args.Bool(2)
}

func (m *MockLockfile) SetDeck(id, parentID, path, name string) {
	m.Called(id, parentID, path, name)
}

func (m *MockLockfile) UpdateDeck(id, name string) {
	m.Called(id, name)
}

func (m *MockLockfile) DeleteDeck(id string) {
	m.Called(id)
}

func (m *MockLockfile) DeleteCard(deckID, cardID string) {
	m.Called(deckID, cardID)
}
