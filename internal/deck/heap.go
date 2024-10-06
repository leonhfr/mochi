package deck

import "path/filepath"

// HeapItem represents the path to a directory and its files.
type HeapItem struct {
	base  string
	paths []string
}

// Heap represents a priority queue for directories.
type Heap []HeapItem

func (h Heap) Len() int           { return len(h) }                          // Len implements heap.Interface.
func (h Heap) Less(i, j int) bool { return len(h[i].base) < len(h[j].base) } // Less implements heap.Interface.
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }                // Swap implements heap.Interface.

// Push implements heap.Interface.
func (h *Heap) Push(x any) {
	path := x.(string)
	base := filepath.Dir(path)
	for i, item := range *h {
		if item.base == base {
			(*h)[i].paths = append((*h)[i].paths, path)
			return
		}
	}
	*h = append(*h, HeapItem{base, []string{path}})
}

// Pop implements heap.Interface.
func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
