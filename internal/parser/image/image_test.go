package image

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/mochi"
)

func Test_Map_Add(t *testing.T) {
	path := "/testdata/Markdown.md"
	tests := []struct {
		name        string
		images      Map
		destination string
		altText     string
		want        Map
	}{
		{
			name:        "should not add when destination is an URL",
			images:      New(path),
			destination: "example.com/image.png",
			want:        New(path),
		},
		{
			name: "should not add when already in map",
			images: Map{
				dirPath: "./testdata",
				images: map[string]Image{
					"testdata/scream.png": {filename: "a42069093fdb614a", destination: "./scream.png", extension: "png", mimeType: "image/png"},
				},
			},
			destination: "./scream.png",
			want: Map{
				dirPath: "./testdata",
				images: map[string]Image{
					"testdata/scream.png": {filename: "a42069093fdb614a", destination: "./scream.png", extension: "png", mimeType: "image/png"},
				},
			},
		},
		{
			name:        "should not add when mime type does not match",
			images:      New(path),
			destination: "./markdown.md",
			want:        New(path),
		},
		{
			name:        "should add",
			images:      New(path),
			destination: "./scream.png",
			want: Map{
				dirPath: "./testdata",
				images: map[string]Image{
					"testdata/scream.png": {filename: "a42069093fdb614a", destination: "./scream.png", extension: "png", mimeType: "image/png"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.images.Add(tt.destination, tt.altText)
			assert.Equal(t, tt.want, tt.images)
		})
	}
}

func Test_Map_Replace(t *testing.T) {
	images := Map{
		dirPath: "./testdata",
		images: map[string]Image{
			"testdata/scream.png":         {filename: "scream_hash", destination: "./scream.png", extension: "png", mimeType: "image/png", altText: "Scream"},
			"testdata/constellations.png": {filename: "constellations_hash", destination: "./constellations.jpg", extension: "jpg", mimeType: "image/jpg"},
		},
	}
	source := "![Scream](./scream.png)\n![](./constellations.jpg)"
	want := "![Scream](@media/scream_hash.png)\n![](@media/constellations_hash.jpg)"
	got := images.Replace(source)
	assert.Equal(t, want, got)
}

func Test_Map_Attachments(t *testing.T) {
	images := Map{
		dirPath: "./testdata",
		images: map[string]Image{
			"testdata/scream.png": {filename: "a42069093fdb614a", destination: "./scream.png", extension: "png", mimeType: "image/png"},
		},
	}
	want := []Attachment{
		{
			Mochi: mochi.Attachment{
				FileName:    "a42069093fdb614a.png",
				ContentType: "image/png",
				Data:        "Q09OVEVO",
			},
			Hash: "45685e95985e20822fb2538a522a5ccf",
			Path: "testdata/scream.png",
		},
	}

	r := new(mockReader)
	r.On("Read", "testdata/scream.png").Return("CONTENT", nil)

	got, err := images.Attachments(r)

	assert.Equal(t, want, got)
	assert.NoError(t, err)
	r.AssertExpectations(t)
}

type mockReader struct {
	mock.Mock
}

func (m *mockReader) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}
