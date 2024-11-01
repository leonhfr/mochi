package image

import (
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/internal/parser"
)

func Test_newImage(t *testing.T) {
	path := "/testdata/Markdown.md"
	tests := []struct {
		name   string
		call   testRead
		parsed parser.Image
		want   Image
		ok     bool
	}{
		{
			name: "should return false when error",
			call: testRead{
				path: "/testdata/scream.png",
				err:  fs.ErrNotExist,
			},
			parsed: parser.Image{
				Destination: "scream.png",
			},
		},
		{
			name: "should return image",
			call: testRead{
				path:    "/testdata/scream.png",
				content: "IMAGE CONTENT",
			},
			parsed: parser.Image{
				Destination: "scream.png",
				AltText:     "alt text",
			},
			want: Image{
				Bytes:       []byte("IMAGE CONTENT"),
				Filename:    "22abb8f07c02970e.png",
				destination: "scream.png",
				altText:     "alt text",
			},
			ok: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader([]testRead{tt.call})
			got, ok := newImage(r, path, tt.parsed)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
			r.AssertExpectations(t)
		})
	}
}

func Test_readImage(t *testing.T) {
	tests := []struct {
		name  string
		call  testRead
		path  string
		bytes []byte
		err   error
	}{
		{
			name: "should read image",
			call: testRead{
				path:    "/testdata/scream.png",
				content: "IMAGE CONTENT",
			},
			path:  "/testdata/scream.png",
			bytes: []byte("IMAGE CONTENT"),
		},
		{
			name: "should return error",
			call: testRead{
				path:    "/testdata/scream.png",
				content: "",
				err:     fs.ErrNotExist,
			},
			path: "/testdata/scream.png",
			err:  fs.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader([]testRead{tt.call})
			content, err := readImage(r, tt.path)
			assert.Equal(t, tt.bytes, content)
			assert.Equal(t, tt.err, err)
			r.AssertExpectations(t)
		})
	}
}

type testRead struct {
	path    string
	content string
	err     error
}

type mockFile struct {
	mock.Mock
}

func newMockReader(calls []testRead) *mockFile {
	m := new(mockFile)
	for _, call := range calls {
		m.On("Read", call.path).Return(call.content, call.err)
	}
	return m
}

func (m *mockFile) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}
