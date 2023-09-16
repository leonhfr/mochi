package sync

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/filesystem"
)

var _ filesystem.Interface = &MockFilesystem{}

type MockFilesystem struct {
	mock.Mock
}

func (m *MockFilesystem) FileExists(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockFilesystem) Read(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFilesystem) Write(path, content string) error {
	args := m.Called(path, content)
	return args.Error(0)
}
