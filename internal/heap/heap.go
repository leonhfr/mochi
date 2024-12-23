package heap

import "container/heap"

// Heap represents a priority heap.
type Heap[T Item] struct {
	heap *priorityHeap[T]
}

// New initializes and returns a new Heap.
func New[T Item]() *Heap[T] {
	h := &priorityHeap[T]{}
	heap.Init(h)
	return &Heap[T]{h}
}

// Len returns the heap length.
func (h *Heap[T]) Len() int {
	return h.heap.Len()
}

// Push pushes a new item to the heap.
func (h *Heap[T]) Push(item T) {
	heap.Push(h.heap, item)
}

// Pop returns the heap item with the most priority (lowest).
func (h *Heap[T]) Pop() Group[T] {
	return heap.Pop(h.heap).(Group[T])
}

// Drain returns all the heap items ordered by priority (lowest).
func (h *Heap[T]) Drain() []Group[T] {
	var groups []Group[T]
	for h.Len() > 0 {
		groups = append(groups, h.Pop())
	}
	return groups
}

// Item is the interface that should be implemented
// for an item to be grouped and prioritized.
type Item interface {
	Base() string
	Priority() int
}

// Group contains the a group of items.
type Group[T Item] struct {
	Base     string
	Items    []T
	priority int
}

type priorityHeap[T Item] []Group[T]

// Less implements heap.Interface.
func (h priorityHeap[T]) Less(i, j int) bool {
	if h[i].priority != h[j].priority {
		return h[i].priority < h[j].priority
	}

	return h[i].Base < h[j].Base
}

func (h priorityHeap[T]) Swap(i, j int) { h[i], h[j] = h[j], h[i] } // Swap implements heap.Interface.
func (h priorityHeap[T]) Len() int      { return len(h) }           // Len implements heap.Interface.

// Push implements heap.Interface.
func (h *priorityHeap[T]) Push(x any) {
	newItem := x.(T)
	for i, item := range *h {
		if item.Base == newItem.Base() {
			(*h)[i].Items = append((*h)[i].Items, newItem)
			return
		}
	}
	*h = append(*h, Group[T]{
		Base:     newItem.Base(),
		Items:    []T{newItem},
		priority: newItem.Priority(),
	})
}

// Pop implements heap.Interface.
func (h *priorityHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
