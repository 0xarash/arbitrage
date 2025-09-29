package arbitrage

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/0xarash/arbitrage/src/decimal"
	"github.com/0xarash/arbitrage/src/limiter"
	"golang.org/x/sync/errgroup"
)

type Market struct {
	config       Config
	exchange     Exchange
	requestLimit int
}

func NewMarket(config Config, exchange Exchange) (Market, error) {
	limit, _, err := exchange.RequestLimit()
	if err != nil {
		return Market{}, err
	}
	market := Market{
		config:       config,
		exchange:     exchange,
		requestLimit: limit,
	}
	return market, nil
}

func (m Market) Pairs() (Pairs, error) {
	pairs, _, err := m.exchange.FetchMarkets()
	return pairs, err
}

type Job struct {
	Pair   Pair
	Limit  string
	Depth  *Depth
	Header http.Header
	Err    error
}

type PairDepth struct {
	Pair  Pair
	Depth Depth
}

type Depths []PairDepth

func (m Market) Depths(pairs Pairs, depthLimit string) (
	Depths, error) {

	curWeight, err := m.getCurrentWeight()
	if err != nil {
		return Depths{}, err
	}

	wg, ctx := errgroup.WithContext(context.Background())
	wg.SetLimit(m.config.Worker.Concurrency)

	limiter := limiter.New(m.requestLimit-curWeight, time.Minute,
		m.config.Limiter.WeightDepth)

	results := make(Depths, len(pairs))

	index := -1
	for sym, pair := range pairs {
		symbol := string(sym)
		pair := pair
		index++
		wg.Go(func() error {
			if err := limiter.Wait(ctx); err != nil {
				return err
			}
			depth, _, err := m.exchange.FetchDepth(symbol, depthLimit)
			if err != nil {
				return err
			}

			results[index] = PairDepth{
				Pair:  pair,
				Depth: depth,
			}
			return nil
		})
	}

	wg.Wait()

	return results, nil
}

type KlineResponse struct {
	Pair  Pair
	Kline []Kline
}

type Klines []KlineResponse

func (m Market) Klines(pairs Pairs,
	limit string) (Klines, error) {

	curWeight, err := m.getCurrentWeight()
	if err != nil {
		return Klines{}, err
	}

	wg, ctx := errgroup.WithContext(context.Background())
	wg.SetLimit(m.config.Worker.Concurrency)

	limiter := limiter.New(m.requestLimit-curWeight, time.Minute,
		m.config.Limiter.WeightKline)

	results := make(Klines, len(pairs))

	index := -1
	for sym, pair := range pairs {
		symbol := string(sym)
		pair := pair
		index++
		wg.Go(func() error {
			if err := limiter.Wait(ctx); err != nil {
				return err
			}
			kline, _, err := m.exchange.FetchKline(map[string]string{
				"symbol":   symbol,
				"limit":    limit,
				"interval": "1d",
			})
			if err != nil {
				return err
			}

			results[index] = KlineResponse{
				Pair:  pair,
				Kline: kline,
			}

			return nil
		})
	}

	err = wg.Wait()
	return results, err
}

func (m Market) getCurrentWeight() (int, error) {
	hdr, err := m.exchange.Ping()
	if err != nil {
		return 0, nil
	}
	curWeight, err := strconv.Atoi(hdr.Get("X-Mbx-Used-Weight-1m"))
	if err != nil {
		return 0, err
	}
	return curWeight, nil
}

type Price struct {
	Pair  Pair
	Price decimal.ExtendedDecimal
}

type Prices map[Symbol]*Price

func (m Market) Prices(pairs Pairs) (Prices, error) {
	prices, _, err := m.exchange.FetchPrice()
	if err != nil {
		return nil, err
	}

	result := make(Prices)
	for _, price := range prices {
		if pair, ok := pairs[price.Symbol]; ok {
			result[price.Symbol] = &Price{
				Pair:  pair,
				Price: price.Price,
			}
		}
	}
	return result, nil
}
