package deck

import (
	"path/filepath"
	"strings"
)

// Directory represents the path to a directory and its files.
type Directory struct {
	Base  string
	Paths []string
	level int
}

// Heap represents a priority queue for directories.
type Heap []Directory

func (h Heap) Len() int           { return len(h) }                  // Len implements heap.Interface.
func (h Heap) Less(i, j int) bool { return h[i].level < h[j].level } // Less implements heap.Interface.
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }        // Swap implements heap.Interface.

// Push implements heap.Interface.
func (h *Heap) Push(x any) {
	path := x.(string)
	base := filepath.Dir(path)
	for i, item := range *h {
		if item.Base == base {
			(*h)[i].Paths = append((*h)[i].Paths, path)
			return
		}
	}
	*h = append(*h, Directory{
		Base:  base,
		Paths: []string{path},
		level: getLevel(base),
	})
}

// Pop implements heap.Interface.
func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getLevel(base string) int {
	if base == "/" {
		return 0
	}
	return strings.Count(base, "/")
}
