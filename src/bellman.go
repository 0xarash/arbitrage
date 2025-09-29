package arbitrage

import (
	"github.com/0xarash/arbitrage/src/decimal"
)

type BellManFord[T comparable] struct {
	graph *Graph[T]
	start T
}

func NewBellmanFord[T comparable](g *Graph[T], start T) BellManFord[T] {
	return BellManFord[T]{
		graph: g,
		start: start,
	}
}

func (b BellManFord[T]) DetectCycles() (dist map[T]decimal.ExtendedDecimal,
	prev map[T]T, cycles []T) {

	dist = make(map[T]decimal.ExtendedDecimal)
	prev = make(map[T]T)

	var empty T
	for _, v := range b.graph.vertex.Values() {
		dist[v], _ = decimal.NewFromFloat64(0, decimal.Infinite)
		prev[v] = empty
	}

	dist[b.start] = decimal.Zero

	for range b.graph.vertex.Len() - 1 {
		for _, edge := range b.graph.Edges() {
			s1, _ := dist[edge.From].Add(edge.Weight)
			if s1.Less(dist[edge.To]) {
				dist[edge.To] = s1
				prev[edge.To] = edge.From
			}
		}
	}

	for _, edge := range b.graph.Edges() {
		s1, _ := dist[edge.From].Add(edge.Weight)
		if s1.Less(dist[edge.To]) {
			cycles = append(cycles, edge.To)
		}
	}

	return dist, prev, cycles
}
