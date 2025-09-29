package arbitrage

import (
	"github.com/0xarash/arbitrage/src/decimal"
	"github.com/rs/zerolog/log"
)

type TriangularArbitrage struct {
}

func NewTriangularArbitrage() TriangularArbitrage {
	return TriangularArbitrage{}
}

func (a TriangularArbitrage) Graph(depths Depths,
	fee decimal.ExtendedDecimal) Graph[Currency] {

	graph := NewGraph[Currency]()
	for _, pair := range depths {
		if pair == nil {
			log.Warn().Msg("Nil Depth")
			continue
		}
		graph.AddVertex(pair.Pair.Base)
		graph.AddVertex(pair.Pair.Quote)

		effectivePrice, _ := decimal.One.Sub(fee)

		price, _ := pair.Depth.Bids[0].Price.Mul(effectivePrice)
		w, _ := price.Log()
		wBaseToQuote := w.Neg()
		graph.AddWeightedEdge(pair.Pair.Base, pair.Pair.Quote, wBaseToQuote)

		price, _ = pair.Depth.Asks[0].Price.Quo(effectivePrice)
		w, _ = price.Log()
		wQuoteToBase := w.Neg()
		graph.AddWeightedEdge(pair.Pair.Quote, pair.Pair.Base, wQuoteToBase)
	}
	return graph
}

func (a TriangularArbitrage) Run(graph Graph[Currency], start Currency) (
	[]Currency, error) {

	bell := NewBellmanFord(&graph, start)
	_, prev, _ := bell.DetectCycles()

	v := start
	for range graph.Vertex().Values() {
		v = prev[v]
	}

	visited := make(map[Currency]bool)
	cycle := []Currency{}
	for !visited[v] {
		visited[v] = true
		cycle = append(cycle, v)
		v = prev[v]
	}

	return cycle, nil
}
