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
		{Base: "/", paths: []string{"/README.md"}},
		{Base: "/lorem-ipsum", paths: []string{"/lorem-ipsum/Lorem ipsum.md", "/lorem-ipsum/Sed interdum libero.md"}},
		{Base: "/lorem-ipsum/Notes", paths: []string{"/lorem-ipsum/Notes/Note 1.md", "/lorem-ipsum/Notes/Note 2.md"}},
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
