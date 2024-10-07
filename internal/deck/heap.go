package deck

import (
	"container/heap"
	"path/filepath"
	"strings"
)

// Directory represents the path to a directory and its files.
type Directory struct {
	Path      string
	FilePaths []string
	level     int
}

// DirHeap represents a heap of directories.
// Priority is given to lower levels (closer to root).
type DirHeap struct {
	dirs *dirHeap
}

// NewDirHeap initializes and returns a new DirHeap.
func NewDirHeap() *DirHeap {
	h := &dirHeap{}
	heap.Init(h)
	return &DirHeap{h}
}

// Len returns the heap length.
func (h *DirHeap) Len() int {
	return h.dirs.Len()
}

// Push pushes a new path to the heap.
func (h *DirHeap) Push(path string) {
	heap.Push(h.dirs, path)
}

// Pop returns the heap directory closest to the root.
func (h *DirHeap) Pop() Directory {
	return heap.Pop(h.dirs).(Directory)
}

type dirHeap []Directory

func (h dirHeap) Len() int           { return len(h) }                  // Len implements heap.Interface.
func (h dirHeap) Less(i, j int) bool { return h[i].level < h[j].level } // Less implements heap.Interface.
func (h dirHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }        // Swap implements heap.Interface.

// Push implements heap.Interface.
func (h *dirHeap) Push(x any) {
	filePath := x.(string)
	path := filepath.Dir(filePath)
	for i, item := range *h {
		if item.Path == path {
			(*h)[i].FilePaths = append((*h)[i].FilePaths, filePath)
			return
		}
	}
	*h = append(*h, Directory{
		Path:      path,
		FilePaths: []string{filePath},
		level:     getLevel(path),
	})
}

// Pop implements heap.Interface.
func (h *dirHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getLevel(path string) int {
	if path == "/" {
		return 0
	}
	return strings.Count(path, "/")
}
