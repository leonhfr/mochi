package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/config"
)

type Config struct {
	Deck []ConfigDeck
}

type ConfigDeck struct {
	Path string
	Deck config.Deck
	OK   bool
}

func NewMockConfig(calls Config) *MockConfig {
	cfg := new(MockConfig)
	for _, call := range calls.Deck {
		cfg.
			On("Deck", call.Path).
			Return(call.Deck, call.OK)
	}
	return cfg
}

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) Deck(path string) (config.Deck, bool) {
	args := m.Called(path)
	return args.Get(0).(config.Deck), args.Bool(1)
}
