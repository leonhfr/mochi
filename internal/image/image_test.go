package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Map_Add(t *testing.T) {
	path := "/testdata/Markdown.md"
	tests := []struct {
		name        string
		images      map[string]Image
		calls       map[string]bool
		destination string
		altText     string
		want        map[string]Image
	}{
		{
			name:        "should not add when destination is an URL",
			calls:       map[string]bool{"testdata/example.com/image.png": false},
			destination: "example.com/image.png",
			want:        map[string]Image{},
		},
		{
			name: "should not add when already in map",
			images: map[string]Image{
				"testdata/scream.png": {Filename: "a42069093fdb614a", Destination: "./scream.png", Extension: "png", MimeType: "image/png"},
			},
			calls:       map[string]bool{"testdata/scream.png": true},
			destination: "./scream.png",
			want: map[string]Image{
				"testdata/scream.png": {Filename: "a42069093fdb614a", Destination: "./scream.png", Extension: "png", MimeType: "image/png"},
			},
		},
		{
			name:        "should not add when mime type does not match",
			calls:       map[string]bool{"testdata/markdown.md": true},
			destination: "./markdown.md",
			want:        map[string]Image{},
		},
		{
			name:        "should add",
			calls:       map[string]bool{"testdata/scream.png": true},
			destination: "./scream.png",
			want: map[string]Image{
				"testdata/scream.png": {Filename: "a42069093fdb614a", Destination: "./scream.png", Extension: "png", MimeType: "image/png"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newMockFileChecker(tt.calls)
			imageMap := New(fc, path)
			imageMap.Add(tt.destination, tt.altText)
			assert.Equal(t, tt.want, imageMap.images)
			fc.AssertExpectations(t)
		})
	}
}

func Test_Map_Replace(t *testing.T) {
	images := Map{
		dirPath: "./testdata",
		images: map[string]Image{
			"testdata/scream.png":         {Filename: "scream_hash", Destination: "./scream.png", Extension: "png", MimeType: "image/png", AltText: "Scream"},
			"testdata/constellations.png": {Filename: "constellations_hash", Destination: "./constellations.jpg", Extension: "jpg", MimeType: "image/jpg"},
		},
	}
	source := "![Scream](./scream.png)\n![](./constellations.jpg)"
	want := "![Scream](@media/scream_hash.png)\n![](@media/constellations_hash.jpg)"
	got := images.Replace(source)
	assert.Equal(t, want, got)
}

type mockFileChecker struct {
	mock.Mock
}

func newMockFileChecker(calls map[string]bool) *mockFileChecker {
	m := new(mockFileChecker)
	for path, ok := range calls {
		m.On("Exists", path).Return(ok)
	}
	return m
}

func (m *mockFileChecker) Exists(p string) bool {
	args := m.Mock.Called(p)
	return args.Bool(0)
}
