package pqueue_test

import(
	"sync"
	"testing"

	"github.com/lkorsman/pqueue"
)

// --- Min-heap tests ---

func TestMinHeap_BasicOrder(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	pq.PushAll(5, 1, 3, 2, 4)

	expected := []int{1, 2, 3, 4, 5}
	for _, want := range expected {
		got, ok := pq.Pop()
		if !ok {
			t.Fatal("Pop returned false unexpectedly")
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestMaxHeap_BasicOrder(t *testing.T) {
	pq := pqueue.New[int](pqueue.Max[int]())
	pq.PushAll(5, 1, 3, 2, 4)
 
	expected := []int{5, 4, 3, 2, 1}
	for _, want := range expected {
		got, ok := pq.Pop()
		if !ok {
			t.Fatal("Pop returned false unexpectedly")
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestPeek_DoesNotRemove(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	pq.PushAll(3, 1, 2)
 
	for i := 0; i < 3; i++ {
		val, ok := pq.Peek()
		if !ok || val != 1 {
			t.Errorf("Peek() = (%d, %v), want (1, true)", val, ok)
		}
	}
	if pq.Len() != 3 {
		t.Errorf("Len() = %d after Peek, want 3", pq.Len())
	}
}
 
func TestPop_EmptyQueue(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	val, ok := pq.Pop()
	if ok {
		t.Errorf("Pop on empty queue returned ok=true, val=%d", val)
	}
	if val != 0 {
		t.Errorf("Pop on empty queue returned non-zero val=%d", val)
	}
}

func TestPeek_EmptyQueue(t *testing.T) {
	pq := pqueue.New[string](pqueue.Min[string]())
	val, ok := pq.Peek()
	if ok || val != "" {
		t.Errorf("Peek on empty queue returned (%q, %v), want (\"\", false)", val, ok)
	}
}
 
func TestIsEmpty(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	if !pq.IsEmpty() {
		t.Error("new queue should be empty")
	}
	pq.Push(1)
	if pq.IsEmpty() {
		t.Error("queue with element should not be empty")
	}
	pq.Pop()
	if !pq.IsEmpty() {
		t.Error("queue after draining should be empty")
	}
}
 
func TestLen(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	for i := 1; i <= 5; i++ {
		pq.Push(i)
		if pq.Len() != i {
			t.Errorf("Len() = %d after %d pushes, want %d", pq.Len(), i, i)
		}
	}
}
 
func TestDrain(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	pq.PushAll(3, 1, 4, 1, 5, 9, 2, 6)
 
	drained := pq.Drain()
	if !pq.IsEmpty() {
		t.Error("queue should be empty after Drain")
	}
 
	for i := 1; i < len(drained); i++ {
		if drained[i] < drained[i-1] {
			t.Errorf("Drain not in order: drained[%d]=%d < drained[%d]=%d",
				i, drained[i], i-1, drained[i-1])
		}
	}
}

func TestSingleElement(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	pq.Push(42)
 
	val, ok := pq.Peek()
	if !ok || val != 42 {
		t.Errorf("Peek() = (%d, %v), want (42, true)", val, ok)
	}
 
	val, ok = pq.Pop()
	if !ok || val != 42 {
		t.Errorf("Pop() = (%d, %v), want (42, true)", val, ok)
	}
 
	if !pq.IsEmpty() {
		t.Error("queue should be empty after popping sole element")
	}
}
 
func TestDuplicates(t *testing.T) {
	pq := pqueue.New[int](pqueue.Min[int]())
	pq.PushAll(3, 3, 1, 1, 2, 2)
 
	prev, _ := pq.Pop()
	for !pq.IsEmpty() {
		curr, _ := pq.Pop()
		if curr < prev {
			t.Errorf("out of order: got %d after %d", curr, prev)
		}
		prev = curr
	}
}

// --- Custom struct test ---

type Task struct {
	Name     string
	Priority int
}

func TestCustomStruct(t *testing.T) {
	pq := pqueue.New[Task](func(a, b Task) bool {
		return a.Priority < b.Priority
	})

	pq.Push(Task{"low", 10})
	pq.Push(Task{"critical", 1})
	pq.Push(Task{"medium", 5})

	first, _ := pq.Pop()
	if first.Name != "critical" {
		t.Errorf("expected 'critical' task first, got '%s'", first.Name)
	}

	second, _ := pq.Pop()
	if second.Name != "medium" {
		t.Errorf("expected 'medium' task second, got '%s'", second.Name)
	}
}

// --- String heap test ---

func TestStringMinHeap(t *testing.T) {
	pq := pqueue.New[string](pqueue.Min[string]())
	pq.PushAll("banana", "apple", "cherry")

	first, _ := pq.Pop()
	if first != "apple" {
		t.Errorf("got %q, want %q", first, "apple")
	}
}

// --- NewWithCapacity test ---

func TestNewWithCapacity(t *testing.T) {
	pq := pqueue.NewWithCapacity[int](pqueue.Min[int](), 100)
	pq.PushAll(3, 1, 2)
	val, _ := pq.Pop()
	if val != 1 {
		t.Errorf("got %d, want %d", val, 1)
	}
}

// --- SyncPriorityQueue tests ---

func TestSyncPriorityQueue_ConcurrentPush(t *testing.T) {
	pq := pqueue.NewSync[int](pqueue.Min[int]())
	var wg sync.WaitGroup
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			pq.Push(val)
		}(i)
	}
	wg.Wait()

	if pq.Len() != n {
		t.Errorf("Len() = %d after %d concurrent pushes, want %d", pq.Len(), n, n)
	}
}

func TestSyncPriorityQueue_ConcurrentPopOrder(t *testing.T) {
	pq := pqueue.NewSync[int](pqueue.Min[int]())
	pq.PushAll(5, 2, 8, 1, 9, 3)

	// Drain sequentially and verify order
	results := pq.Drain()
	for i := 1; i < len(results); i++ {
		if results[i] < results[i-1] {
			t.Errorf("out of order after drain: %v", results)
		}
	}
}

// --- Benchmarks ---

func BenchmarchPush(b *testing.B) {
	pq := pqueue.NewWithCapacity[int](pqueue.Min[int](), b.N)
	for i := 0; i < b.N; i++ {
		pq.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	pq := pqueue.NewWithCapacity[int](pqueue.Min[int](), b.N)
	for i := 0; i < b.N; i++ {
		pq.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Pop()
	}
}

func BenchmarkPushPop(b *testing.B) {
	pq := pqueue.New[int](pqueue.Min[int]())
	for i := 0; i < b.N; i++ {
		pq.Push(i)
		pq.Pop()
	}
}
