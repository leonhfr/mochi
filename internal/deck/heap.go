package deck

import (
	"path/filepath"
	"strings"
)

// Directory represents the path to a directory and its files.
type Directory struct {
	Path      string
	FilePaths []string
	level     int
}

// Heap represents a priority queue for directories.
type Heap []Directory

func (h Heap) Len() int           { return len(h) }                  // Len implements heap.Interface.
func (h Heap) Less(i, j int) bool { return h[i].level < h[j].level } // Less implements heap.Interface.
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }        // Swap implements heap.Interface.

// Push implements heap.Interface.
func (h *Heap) Push(x any) {
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
func (h *Heap) Pop() any {
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
