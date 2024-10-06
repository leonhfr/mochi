package deck

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Heap(t *testing.T) {
	input := []string{
		"/lorem-ipsum/Lorem ipsum.md",
		"/lorem-ipsum/Notes/Note 1.md",
		"/README.md",
		"/lorem-ipsum/Notes/Note 2.md",
		"/lorem-ipsum/Sed interdum libero.md",
	}
	want := []Directory{
		{level: 0, Path: "/", FilePaths: []string{"/README.md"}},
		{level: 1, Path: "/lorem-ipsum", FilePaths: []string{"/lorem-ipsum/Lorem ipsum.md", "/lorem-ipsum/Sed interdum libero.md"}},
		{level: 2, Path: "/lorem-ipsum/Notes", FilePaths: []string{"/lorem-ipsum/Notes/Note 1.md", "/lorem-ipsum/Notes/Note 2.md"}},
	}

	h := &Heap{}
	heap.Init(h)

	for _, path := range input {
		heap.Push(h, path)
	}

	var got []Directory
	for h.Len() > 0 {
		got = append(got, heap.Pop(h).(Directory))
	}

	assert.Equal(t, want, got)
}
