package test

import (
	"io"
	"strings"

	"github.com/stretchr/testify/mock"
)

type MockFile struct {
	mock.Mock
}

func (m *MockFile) Exists(p string) bool {
	args := m.Mock.Called(p)
	return args.Bool(0)
}

func (m *MockFile) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}

func (m *MockFile) Write(p string) (io.WriteCloser, error) {
	args := m.Mock.Called(p)
	wc := writeCloser{&strings.Builder{}}
	return wc, args.Error(1)
}

type writeCloser struct {
	*strings.Builder
}

func (writeCloser) Close() error { return nil }
