// Package pqueue provides a generic, type-safe priority queue (min-heap or max-heap).
// It works with any type via a user-provided comparison function, and ships with
// built-in Min/Max helpers for ordered types.
package pqueue

import "cmp"

// PriorityQueue is a generic heap-backed priority queue.
// The ordering is determined by the comparator passed to New:
// if less(a, b) is true, a will be popped before b.
type PriorityQueue[T any] struct {
	data []T
	less func(a, b T) bool
}

// New creates a new PriorityQueue ordered by the given comparator.
// If less(a, b) returns true, a has higher priority than b.
//
// Example - min-heap of ints:
//
// pq := pqueue.New[int](pqueue.Min[int]())
//
// Example - max-heap of ints:
//
// pq := pqueue.New[int](pqueue.Max[int]())
//
// # Example - custom struct ordered by a field
//
//	pq := pqueue.New[Task](func(a, b Task) bool {
//			return a.Priority < b.Priority
//	})
func New[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data: []T{},
		less: less,
	}
}

// NewWithCapacity creates aa PriorityQueue with a pre-allocated capacity,
// useful when you know roughly how many elements you'll store.
func NewWithCapacity[T any](less func(a, b T) bool, capacity int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data: make([]T, 0, capacity),
		less: less,
	}
}

// Min returns a comparator that produces a min-heap (smallest value popped first).
// T must be an ordered type (int, float64, string, etc.).
func Min[T cmp.Ordered]() func(a, b T) bool {
	return func(a, b T) bool { return a < b }
}

// Max returns a comparator that produces a max-heap (largest value popped first).
// T must be an ordered type (int, float64, string, etc.).
func Max[T cmp.Ordered]() func(a, b T) bool {
	return func(a, b T) bool { return a > b }
}

// Push adds a value to the priority queue. O(log n).
func (pq *PriorityQueue[T]) Push(val T) {
	pq.data = append(pq.data, val)
	pq.siftUp(len(pq.data) - 1)
}

// Pop removes and returns the highest-priority value.
// Returns the zero value of T and false if the queue is empty. O(log n).
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	if len(pq.data) == 0 {
		var zero T
		return zero, false
	}

	top := pq.data[0]
	last := len(pq.data) - 1

	pq.data[0] = pq.data[last]
	pq.data = pq.data[:last]

	if len(pq.data) > 0 {
		pq.siftDown(0)
	}

	return top, true
}

// Peek returns the highest-priority value without removing it.
// Returns the zero value of T and false if the queue is empty. O(1).
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	if len(pq.data) == 0 {
		var zero T
		return zero, false
	}
	return pq.data[0], true
}

// Len returns the number of elements in the queue. O(1).
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.data)
}

// IsEmpty reports whether the queue has no elements. O(1).
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return len(pq.data) == 0
}

// PushAll adds multiple values to the queue.
func (pq *PriorityQueue[T]) PushAll(vals ...T) {
	for _, v := range vals {
		pq.Push(v)
	}
}

// Drain removes and returns all the elements in priority order.
// The queue is empty after this call.
func (pq *PriorityQueue[T]) Drain() []T {
	result := make([]T, 0, len(pq.data))
	for !pq.IsEmpty() {
		val, _ := pq.Pop()
		result = append(result, val)
	}
	return result
}

// --- internal heap operations ---

func (pq *PriorityQueue[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if pq.less(pq.data[i], pq.data[parent]) {
			pq.data[i], pq.data[parent] = pq.data[parent], pq.data[i]
			i = parent
		} else {
			break
		}
	}
}

func (pq *PriorityQueue[T]) siftDown(i int) {
	n := len(pq.data)
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2

		if left < n && pq.less(pq.data[left], pq.data[right]) {
			smallest = left
		}
		if right < n && pq.less(pq.data[right], pq.data[left]) {
			smallest = right
		}
		if smallest == i {
			break
		}

		pq.data[i], pq.data[smallest] = pq.data[smallest], pq.data[i]
		i = smallest
	}
}
