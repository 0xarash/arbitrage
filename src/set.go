package arbitrage

type Set[T comparable] struct {
	m map[T]any
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]any)}
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) Add(v T) {
	s.m[v] = struct{}{}
}

func (s *Set[T]) Values() []T {
	vals := make([]T, 0, len(s.m))
	for k := range s.m {
		vals = append(vals, k)
	}
	return vals
}
