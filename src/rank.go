package arbitrage

import "github.com/0xarash/arbitrage/src/decimal"

type Ranker struct {
	metrics Metrics
}

func NewRanker(metrics Metrics) Ranker {
	return Ranker{
		metrics: metrics,
	}
}

func (r Ranker) GetTop(count int, priorityThreshold decimal.ExtendedDecimal,
	ignoreZero bool) (Pairs, error) {

	pq := NewPriorityQueue[*Pair]()

	for pair, metric := range r.metrics {
		if ignoreZero {
			if metric.Volume.Amount.IsZero() {
				continue
			}
		}

		if !priorityThreshold.IsZero() {
			if metric.Volume.Amount.Less(priorityThreshold) {
				continue
			}
		}
		pq.Push(&pair, metric.Volume.Amount)
	}

	ranked := make(Pairs)
	for count > 0 && pq.Len() > 0 {
		pair := pq.Pop()
		ranked[pair.Symbol()] = *pair
		count--
	}
	return ranked, nil
}
