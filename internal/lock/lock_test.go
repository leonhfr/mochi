package lock

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/v3/assert"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name        string
		target      string
		path        string
		exists      bool
		fileContent string
		fileError   error
		wantPath    string
		wantData    lockData
		err         error
	}{
		{
			name:     "no lockfile found",
			target:   "testdata",
			path:     "testdata/mochi-lock.json",
			wantPath: "testdata/mochi-lock.json",
		},
		{
			name:        "lockfile found",
			target:      "testdata",
			path:        "testdata/mochi-lock.json",
			exists:      true,
			fileContent: `{"DECK_ID":{"path":"DECK_PATH","name":"DECK_NAME"}}`,
			wantPath:    "testdata/mochi-lock.json",
			wantData:    lockData{"DECK_ID": lockDeck{Path: "DECK_PATH", Name: "DECK_NAME"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := new(mockReadWriter)
			rw.On("Exists", tt.path).Return(tt.exists)
			if tt.exists {
				rw.On("Read", tt.path).Return(tt.fileContent, tt.fileError)
			}

			got, err := Parse(rw, tt.target)
			assert.Equal(t, tt.wantPath, got.path)
			assert.DeepEqual(t, tt.wantData, got.data)
			assert.Equal(t, tt.err, err)
			rw.AssertExpectations(t)
		})
	}
}

type mockReadWriter struct {
	mock.Mock
}

func (m *mockReadWriter) Exists(p string) bool {
	args := m.Mock.Called(p)
	return args.Bool(0)
}

func (m *mockReadWriter) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}

func (m *mockReadWriter) Write(p string) (io.WriteCloser, error) {
	args := m.Mock.Called(p)
	wc := writeCloser{&strings.Builder{}}
	return wc, args.Error(1)
}

type writeCloser struct {
	*strings.Builder
}

func (writeCloser) Close() error { return nil }
