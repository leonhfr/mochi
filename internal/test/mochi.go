package test

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/mochi"
)

var ErrMochi = errors.New("mochi error")

type Mochi struct {
	CreateDeck []MochiCreateDeck
	UpdateDeck []MochiUpdateDeck
}

type MochiCreateDeck struct {
	Req  mochi.CreateDeckRequest
	Deck mochi.Deck
	Err  error
}

type MochiUpdateDeck struct {
	ID   string
	Req  mochi.UpdateDeckRequest
	Deck mochi.Deck
	Err  error
}

type MockMochi struct {
	mock.Mock
}

func NewMockMochi(calls Mochi) *MockMochi {
	m := new(MockMochi)
	for _, call := range calls.CreateDeck {
		m.
			On("CreateDeck", mock.Anything, call.Req).
			Return(call.Deck, call.Err)
	}
	for _, call := range calls.UpdateDeck {
		m.
			On("UpdateDeck", mock.Anything, call.ID, call.Req).
			Return(call.Deck, call.Err)
	}
	return m
}

func (m *MockMochi) CreateDeck(ctx context.Context, req mochi.CreateDeckRequest) (mochi.Deck, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(mochi.Deck), args.Error(1)
}

func (m *MockMochi) UpdateDeck(ctx context.Context, id string, req mochi.UpdateDeckRequest) (mochi.Deck, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(mochi.Deck), args.Error(1)
}
