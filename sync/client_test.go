package sync

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/api"
)

var (
	_ Client = &api.Client{}
	_ Client = &MockClient{}
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) ListDecks(ctx context.Context) ([]api.Deck, error) {
	args := m.Called(ctx)
	return args.Get(0).([]api.Deck), args.Error(1)
}

func (m *MockClient) ListTemplates(ctx context.Context) ([]api.Template, error) {
	args := m.Called(ctx)
	return args.Get(0).([]api.Template), args.Error(1)
}
