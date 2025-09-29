package arbitrage

import (
	"github.com/0xarash/arbitrage/src/decimal"
	"github.com/rs/zerolog/log"
)

type Metric struct {
	Volume Money
}

type Metrics map[Pair]Metric

func RankMetrics(pairs Pairs, klines Klines, conv Converter,
	ref Currency) Metrics {

	result := make(Metrics)
	for _, resp := range klines {
		if pair, ok := pairs[resp.Pair.Symbol()]; ok {
			totalVolume := decimal.Zero
			for _, kline := range resp.Kline {
				totalVolume, _ = totalVolume.Add(kline.Volume)
			}
			vol, err := conv.Convert(Money{Currency: pair.Base, Amount: totalVolume},
				ref,
			)
			if err != nil {
				log.Warn().Err(err)
				continue
			}
			result[pair] = Metric{
				Volume: vol,
			}
		}
	}
	return result
}
