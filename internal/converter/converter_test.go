package converter

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Converter_Convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		calls  []testRead
		source string
		want   Result
	}{
		{
			name:   "should convert markdown",
			path:   "/testdata/Markdown.md",
			source: "# Hello, World!\n",
			want:   Result{Markdown: "# Hello, World!\n"},
		},
		{
			name: "images",
			path: "/testdata/Images.md",
			calls: []testRead{
				{path: "/testdata/scream.png", content: "IMAGE CONTENT", err: nil},
			},
			source: "![Scream](./scream.png)\n",
			want: Result{
				Markdown: "![Scream](@media/22abb8f07c02970e.png)\n",
				Attachments: []Attachment{
					{Bytes: []byte("IMAGE CONTENT"), Filename: "22abb8f07c02970e.png"},
				},
			},
		},
		{
			name:   "video",
			path:   "/testdata/Video.md",
			source: "![@video](youtube.com/VIDEO)\n",
			want: Result{
				Markdown: "<iframe src=\"youtube.com/VIDEO?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0\" frameborder=\"0\" loading=\"lazy\" gesture=\"media\" allow=\"autoplay; fullscreen\" allowautoplay=\"true\" allowfullscreen=\"true\" style=\"aspect-ratio:16/9;height:100%;width:100%;\"></iframe>\n",
			},
		},
		{
			name:   "mermaid",
			source: "```mermaid\ngraph TD;\n    A-->B;\n    A-->C;\n    B-->D;\n    C-->D;\n```\n",
			want: Result{
				Markdown: "<pre class=\"mermaid\">graph TD;\n    A--&gt;B;\n    A--&gt;C;\n    B--&gt;D;\n    C--&gt;D;\n</pre><script src=\"https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js\"></script><script>mermaid.initialize({startOnLoad: true});</script>\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader(tt.calls)
			c := New()
			got, err := c.Convert(r, tt.path, tt.source)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
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
