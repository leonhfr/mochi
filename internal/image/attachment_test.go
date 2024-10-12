package image

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/mochi/mochi"
)

func Test_Attachments(t *testing.T) {
	images := map[string]Image{
		"testdata/scream.png": {Filename: "a42069093fdb614a", Extension: "png", MimeType: "image/png"},
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

	got, err := Attachments(r, images)

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
