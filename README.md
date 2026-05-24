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