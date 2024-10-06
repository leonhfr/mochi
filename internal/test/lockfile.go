package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/lock"
)

type Lockfile struct {
	GetDeck        []LockfileGetDeck
	SetDeck        []LockfileSetDeck
	UpdateDeckName []LockfileUpdateDeckName
}

type LockfileGetDeck struct {
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

type LockfileUpdateDeckName struct {
	ID   string
	Name string
}

func NewMockLockfile(calls Lockfile) *MockLockfile {
	lf := new(MockLockfile)
	for _, call := range calls.GetDeck {
		lf.
			On("GetDeck", call.Path).
			Return(call.DeckID, call.Deck, call.OK)
	}
	for _, call := range calls.SetDeck {
		lf.
			On("SetDeck", call.ID, call.ParentID, call.Path, call.Name).
			Return()
	}
	for _, call := range calls.UpdateDeckName {
		lf.
			On("UpdateDeckName", call.ID, call.Name).
			Return()
	}
	return lf
}

type MockLockfile struct {
	mock.Mock
}

func (m *MockLockfile) GetDeck(path string) (string, lock.Deck, bool) {
	args := m.Called(path)
	return args.String(0), args.Get(1).(lock.Deck), args.Bool(2)
}

func (m *MockLockfile) SetDeck(id, parentID, path, name string) {
	m.Called(id, parentID, path, name)
}

func (m *MockLockfile) UpdateDeckName(id, name string) {
	m.Called(id, name)
}
