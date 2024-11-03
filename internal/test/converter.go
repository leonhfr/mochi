package test

import (
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/converter"
)

type ConverterCall struct {
	Path   string
	Source string
	Result converter.Result
	Err    error
}

type MockConverter struct {
	mock.Mock
}

func NewMockConverter(calls []ConverterCall) *MockConverter {
	m := new(MockConverter)
	for _, call := range calls {
		m.
			On("Convert", mock.Anything, call.Path, call.Source).
			Return(call.Result, call.Err)
	}
	return m
}

func (m *MockConverter) Convert(reader converter.Reader, path, source string) (converter.Result, error) {
	args := m.Called(reader, path, source)
	return args.Get(0).(converter.Result), args.Error(1)
}
