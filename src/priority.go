package arbitrage

import (
	"container/heap"

	"github.com/0xarash/arbitrage/src/decimal"
)

type PriorityQueue[T any] struct {
	internal _priorityQueue[T]
}

func NewPriorityQueue[T any]() PriorityQueue[T] {
	pq := PriorityQueue[T]{}
	return pq
}

func (pq PriorityQueue[T]) Len() int { return pq.internal.Len() }

func (pq *PriorityQueue[T]) Push(value T, priority decimal.ExtendedDecimal) {
	heap.Push(&pq.internal, &priorityItem[T]{
		value:    value,
		priority: priority,
	})
}

func (pq *PriorityQueue[T]) Pop() T {
	item := heap.Pop(&pq.internal).(*priorityItem[T])
	return item.value
}

type priorityItem[T any] struct {
	value    T
	priority decimal.ExtendedDecimal
}
type _priorityQueue[T any] struct {
	items []*priorityItem[T]
}

func (pq _priorityQueue[T]) Len() int { return len(pq.items) }

func (pq _priorityQueue[T]) Less(i, j int) bool {
	return !pq.items[i].priority.Less(pq.items[j].priority)
}

func (pq _priorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *_priorityQueue[T]) Push(x any) {
	item := x.(*priorityItem[T])
	pq.items = append(pq.items, item)
}

func (pq *_priorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	pq.items = old[0 : n-1]
	return item
}
