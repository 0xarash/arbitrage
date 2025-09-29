package arbitrage

import "github.com/0xarash/arbitrage/src/decimal"

type Edge[T comparable] struct {
	From   T
	To     T
	Weight decimal.ExtendedDecimal
}

type Graph[T comparable] struct {
	vertex *Set[T]
	edges  []*Edge[T]
}

func NewGraph[T comparable]() Graph[T] {
	return Graph[T]{
		vertex: NewSet[T](),
	}
}

func (g *Graph[T]) AddWeightedEdge(from, to T, w decimal.ExtendedDecimal) {
	g.edges = append(g.edges, &Edge[T]{
		From:   from,
		To:     to,
		Weight: w,
	})
}

func (g *Graph[T]) AddVertex(v T) {
	g.vertex.Add(v)
}

func (g *Graph[T]) Edges() []*Edge[T] {
	return g.edges
}

func (g *Graph[T]) Vertex() *Set[T] {
	return g.vertex
}
