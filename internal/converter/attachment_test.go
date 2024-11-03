package converter

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newAttachment(t *testing.T) {
	path := "/testdata/Markdown.md"
	tests := []struct {
		name        string
		call        testRead
		destination string
		want        Attachment
		err         error
	}{
		{
			name: "should return false when error",
			call: testRead{
				path: "/testdata/scream.png",
				err:  fs.ErrNotExist,
			},
			destination: "scream.png",
			err:         fs.ErrNotExist,
		},
		{
			name: "should return image",
			call: testRead{
				path:    "/testdata/scream.png",
				content: "IMAGE CONTENT",
			},
			destination: "scream.png",
			want: Attachment{
				Bytes:    []byte("IMAGE CONTENT"),
				Filename: "22abb8f07c02970e.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader([]testRead{tt.call})
			got, err := newAttachment(r, path, tt.destination)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
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
			content, err := readAttachment(r, tt.path)
			assert.Equal(t, tt.bytes, content)
			assert.Equal(t, tt.err, err)
			r.AssertExpectations(t)
		})
	}
}
