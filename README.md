# pqueue

A generic, type-safe priority queue for Go — backed by a binary heap with zero dependencies.

Go's standard library includes `container/heap`, but using it requires implementing a 5-method interface for every type you want to prioritize. `pqueue` eliminates that boilerplate with generics and a simple comparator function.

## Requirements

Go 1.21+

## Installation

```bash
go get github.com/lkorsman/pqueue
```

## Quick Start

```go
// Min-heap: smallest value comes out first
pq := pqueue.New[int](pqueue.Min[int]())
pq.PushAll(5, 1, 3, 2, 4)

for !pq.IsEmpty() {
    val, _ := pq.Pop()
    fmt.Println(val) // 1, 2, 3, 4, 5
}
```

```go
// Max-heap: largest value comes out first
pq := pqueue.New[int](pqueue.Max[int]())
pq.PushAll(5, 1, 3, 2, 4)

val, _ := pq.Pop() // 5
```

## Custom Types

Pass your own comparator to order any struct:

```go
type Task struct {
    Name     string
    Priority int
}

pq := pqueue.New[Task](func(a, b Task) bool {
    return a.Priority < b.Priority // lower number = higher priority
})

pq.Push(Task{"send email",   10})
pq.Push(Task{"fix outage",    1})
pq.Push(Task{"write tests",   5})

task, _ := pq.Pop()
fmt.Println(task.Name) // "fix outage"
```

## Thread Safety

Use `NewSync` for concurrent access:

```go
pq := pqueue.NewSync[int](pqueue.Min[int]())

// Safe to Push/Pop from multiple goroutines simultaneously
go func() { pq.Push(42) }()
go func() { pq.Push(7)  }()
```

## API Reference

| Method | Description | Complexity |
|---|---|---|
| `Push(val T)` | Add a value | O(log n) |
| `Pop() (T, bool)` | Remove and return highest-priority value | O(log n) |
| `Peek() (T, bool)` | Return highest-priority value without removing | O(1) |
| `Len() int` | Number of elements | O(1) |
| `IsEmpty() bool` | Whether the queue is empty | O(1) |
| `PushAll(vals ...T)` | Add multiple values | O(n log n) |
| `Drain() []T` | Remove and return all values in priority order | O(n log n) |

`SyncPriorityQueue` exposes the same API, safe for concurrent use.

## Design Notes

- **Comparator-based** — works with any type, no interface implementation required
- **No dependencies** — only uses the Go standard library
- **Generics-first** — written for Go 1.21+, fully type-safe
- **Two built-in helpers** — `pqueue.Min[T]()` and `pqueue.Max[T]()` for ordered types

## License

MIT