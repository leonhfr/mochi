package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Heap_Path(t *testing.T) {
	input := []Path{
		"/lorem-ipsum/Lorem ipsum.md",
		"/lorem-ipsum/Notes/Note 1.md",
		"/README.md",
		"/lorem-ipsum/Notes/Note 2.md",
		"/lorem-ipsum/Sed interdum libero.md",
	}
	want := []Group[Path]{
		{priority: 0, Base: "/", Items: []Path{"/README.md"}},
		{priority: 1, Base: "/lorem-ipsum", Items: []Path{"/lorem-ipsum/Lorem ipsum.md", "/lorem-ipsum/Sed interdum libero.md"}},
		{priority: 2, Base: "/lorem-ipsum/Notes", Items: []Path{"/lorem-ipsum/Notes/Note 1.md", "/lorem-ipsum/Notes/Note 2.md"}},
	}

	h := New[Path]()
	for _, path := range input {
		h.Push(path)
	}

	var got []Group[Path]
	for h.Len() > 0 {
		got = append(got, h.Pop())
	}

	assert.Equal(t, want, got)
}
