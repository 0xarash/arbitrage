package arbitrage

import (
	"net/http"
	"time"

	"github.com/0xarash/arbitrage/src/decimal"
)

type Exchange interface {
	Ping() (http.Header, error)
	RequestLimit() (int, http.Header, error)
	FetchMarkets() (Pairs, http.Header, error)
	FetchDepth(symbol, limit string) (*Depth, http.Header, error)
	FetchKline(params map[string]string) ([]Kline, http.Header, error)
	FetchPrice() ([]PriceItem, http.Header, error)
}

type Depth struct {
	Bids []Level `json:"bids"`
	Asks []Level `json:"asks"`
}

type Level struct {
	Price    decimal.ExtendedDecimal
	Quantity decimal.ExtendedDecimal
}

type Kline struct {
	OpenTime  time.Time
	Open      decimal.ExtendedDecimal
	High      decimal.ExtendedDecimal
	Low       decimal.ExtendedDecimal
	Close     decimal.ExtendedDecimal
	Volume    decimal.ExtendedDecimal
	CloseTime time.Time
	Trades    float64
}

type PriceItem struct {
	Symbol Symbol                  `json:"symbol"`
	Price  decimal.ExtendedDecimal `json:"price"`
}
