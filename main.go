package main

import (
	arbitrage "github.com/0xarash/arbitrage/src"
	"github.com/0xarash/arbitrage/src/decimal"
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

func ComputeRankedPairs(cfg arbitrage.Config, market arbitrage.Market) (
	arbitrage.Pairs, error) {

	log.Info().Msg("Fetching Pairs")
	pairs, err := market.Pairs()
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Fetching Prices")
	prices, err := market.Prices(pairs)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("New Converter")
	baseCurrencies := []arbitrage.Currency{arbitrage.USDT}
	converter, err := arbitrage.NewConverter(prices, baseCurrencies)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Fetching Klines")
	klines, err := market.Klines(pairs, "2")
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Calculating Rank Metrics")
	metrics := arbitrage.RankMetrics(pairs, klines, converter, arbitrage.USDT)

	log.Info().Msg("New Ranker")
	ranker := arbitrage.NewRanker(metrics)

	log.Info().Msg("Top Pairs")
	threshold, _ := decimal.New(cfg.Ranking.VolumeThreshold, 0)
	rankedPairs, err := ranker.GetTop(cfg.Ranking.TopPairs, threshold,
		cfg.Ranking.IgnoreZero)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Top Pairs Count = %d", len(rankedPairs))
	for _, pair := range rankedPairs {
		log.Info().Msgf("Symbol: %s, Volume (USDT): %s", pair.Symbol(),
			metrics[pair].Volume.Amount.String())
	}

	return rankedPairs, nil
}

func main() {
	var cfg arbitrage.Config
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	exchange := arbitrage.NewBinance()
	market, err := arbitrage.NewMarket(cfg, exchange)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	rankedPairs, err := ComputeRankedPairs(cfg, market)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msg("Fetching Depth")
	depths, err := market.GetDepths(rankedPairs, "1")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msg("Calculating Arbitrage")
	fee, _ := decimal.NewFromFloat64(cfg.Binance.TradingFee, decimal.Finite)
	triangular := arbitrage.NewTriangularArbitrage()
	graph := triangular.Graph(depths, fee)
	for _, start := range cfg.Arbitrage.StartCurrencies {
		log.Info().Msgf("Start currency: %s", start)
		tradeCycle, err := triangular.Run(graph, start)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		if len(tradeCycle) < 3 {
			continue
		}

		msg := string(tradeCycle[0])
		for i := 1; i < len(tradeCycle); i++ {
			msg += " -> " + string(tradeCycle[i])
		}
		log.Info().Msg(msg)
	}

}
