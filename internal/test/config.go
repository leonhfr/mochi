package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/config"
)

type Config struct {
	GetDeck []ConfigGetDeck
}

type ConfigGetDeck struct {
	Base string
	Deck config.Deck
	OK   bool
}

func NewMockConfig(calls Config) *MockConfig {
	cfg := new(MockConfig)
	for _, call := range calls.GetDeck {
		cfg.
			On("GetDeck", call.Base).
			Return(call.Deck, call.OK)
	}
	return cfg
}

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetDeck(base string) (config.Deck, bool) {
	args := m.Called(base)
	return args.Get(0).(config.Deck), args.Bool(1)
}
