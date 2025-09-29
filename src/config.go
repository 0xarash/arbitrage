package arbitrage

type Config struct {
	Arbitrage struct {
		StartCurrencies []Currency `toml:"start_currencies"`
	}
	Ranking struct {
		VolumeThreshold int64 `toml:"volume_threshold"`
		TopPairs        int   `toml:"top_pairs"`
		IgnoreZero      bool  `toml:"ignore_zero"`
	}
	Worker struct {
		Concurrency int `toml:"concurrency"`
	}
	Limiter struct {
		WeightKline int `toml:"weight_kline"`
		WeightDepth int `toml:"weight_depth"`
	}
	Binance struct {
		TradingFee float64 `toml:"trading_fee"`
	}
}
