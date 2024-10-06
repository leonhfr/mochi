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
	want := []HeapItem{
		{base: "/", paths: []string{"/README.md"}},
		{base: "/lorem-ipsum", paths: []string{"/lorem-ipsum/Lorem ipsum.md", "/lorem-ipsum/Sed interdum libero.md"}},
		{base: "/lorem-ipsum/Notes", paths: []string{"/lorem-ipsum/Notes/Note 1.md", "/lorem-ipsum/Notes/Note 2.md"}},
	}

	h := &Heap{}
	heap.Init(h)

	for _, path := range input {
		heap.Push(h, path)
	}

	var got []HeapItem
	for h.Len() > 0 {
		got = append(got, heap.Pop(h).(HeapItem))
	}

	assert.Equal(t, want, got)
}
