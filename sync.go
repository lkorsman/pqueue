package pqueue
 
import "sync"
 
// SyncPriorityQueue is a thread-safe wrapper around PriorityQueue.
// All operations are safe to call from multiple goroutines concurrently.
type SyncPriorityQueue[T any] struct {
	mu sync.Mutex
	pq *PriorityQueue[T]
}
 
// NewSync creates a new thread-safe priority queue.
//
// Example:
//
//	pq := pqueue.NewSync[int](pqueue.Min[int]())
func NewSync[T any](less func(a, b T) bool) *SyncPriorityQueue[T] {
	return &SyncPriorityQueue[T]{pq: New[T](less)}
}

// Push adds a value to the queue. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) Push(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pq.Push(val)
}

// Pop removes and returns the highest-priority value.
// Returns the zero value and false if the queue is empty. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pq.Pop()
}

// Peek returns the highest-priority value without removing it.
// Returns the zero value and false if the queue is empty. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) Peek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pq.Peek()
}

// Len returns the number of elements in the queue. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pq.Len()
}

// IsEmpty reports whether the queue has no elements. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pq.IsEmpty()
}
 
// PushAll adds multiple values to the queue atomicaly. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) PushAll(vals ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pq.PushAll(vals...)
}

// Drain removes and returns all the elements in priority order.
// The queue is empty after this call. Safe for concurrent use.
func (s *SyncPriorityQueue[T]) Drain() []T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pq.Drain()
}