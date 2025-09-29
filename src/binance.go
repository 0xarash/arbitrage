package arbitrage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/0xarash/arbitrage/src/decimal"
)

type Binance struct{}

func NewBinance() Exchange {
	return Binance{}
}

func (p *Pair) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
		Status     string `json:"status"`
		IsSpot     bool   `json:"isSpotTradingAllowed"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if !tmp.IsSpot || tmp.Status != "TRADING" {
		return nil
	}

	p.Base = Currency(tmp.BaseAsset)
	p.Quote = Currency(tmp.QuoteAsset)
	return nil
}

func (e Binance) FetchMarkets() (Pairs, http.Header, error) {
	body, hdr, err := e.request("https://api.binance.com/api/v3/exchangeInfo",
		map[string]string{})
	if err != nil {
		return nil, nil, errors.New(fmt.Sprint("GetMarkets Failed, ", err))
	}
	defer body.Close()

	d := json.NewDecoder(body)

	empty := false
	for {
		tok, err := d.Token()
		if err != nil {
			return nil, nil, errors.New(fmt.Sprint("GetMarkets Failed, ", err))
		}
		if key, ok := tok.(string); ok && key == "symbols" {
			break
		}
	}

	if empty {
		return nil, nil, errors.New("GetMarkets Failed, empty symbols")
	}

	// omit '['
	d.Token()

	pairs := make(Pairs)
	for d.More() {
		var pair Pair
		if err := d.Decode(&pair); err != nil {
			return nil, nil, errors.New(fmt.Sprint("GetMarkets Failed, ", err))
		}
		if len(pair.Base) == 0 {
			continue
		}
		pairs[pair.Symbol()] = pair
	}

	return pairs, hdr, nil
}

func (l *Level) UnmarshalJSON(data []byte) error {
	var raw [2]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	price, err := decimal.Parse(raw[0])
	if err != nil {

		return errors.New(fmt.Sprint("UnmarshalJSON on Price Failed, ", err))
	}

	qty, err := decimal.Parse(raw[1])
	if err != nil {
		return errors.New(fmt.Sprint("UnmarshalJSON on Qty Failed, ", err))
	}

	l.Price = price
	l.Quantity = qty

	return nil
}

func (e Binance) FetchDepth(symbol, limit string) (
	Depth, http.Header, error) {
	body, hdr, err := e.request("https://api.binance.com/api/v3/depth",
		map[string]string{"symbol": symbol, "limit": limit})
	if err != nil {
		return Depth{}, nil, errors.New(fmt.Sprint("GetDepth Failed, ", err))
	}

	d := json.NewDecoder(body)

	var depth Depth
	if err := d.Decode(&depth); err != nil {
		return Depth{}, nil, errors.New(fmt.Sprint("GetDepth Failed, ", err))
	}

	return depth, hdr, nil
}

func (k *Kline) UnmarshalJSON(data []byte) error {
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if open, ok := raw[1].(string); ok {
		k.Open, _ = decimal.Parse(open)
	}

	if high, ok := raw[2].(string); ok {
		k.High, _ = decimal.Parse(high)
	}

	if low, ok := raw[3].(string); ok {
		k.Low, _ = decimal.Parse(low)
	}

	if close, ok := raw[4].(string); ok {
		k.Close, _ = decimal.Parse(close)
	}

	if vol, ok := raw[5].(string); ok {
		k.Volume, _ = decimal.Parse(vol)
	}

	if trades, ok := raw[8].(float64); ok {
		k.Trades = trades
	}

	k.OpenTime = time.UnixMilli((int64(raw[0].(float64))))
	k.CloseTime = time.UnixMilli((int64(raw[6].(float64))))

	return nil
}

func (e Binance) FetchKline(params map[string]string) (
	[]Kline, http.Header, error) {

	body, hdr, err := e.request("https://api.binance.com/api/v3/klines", params)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprint("GetKline Failed, ", err))
	}

	d := json.NewDecoder(body)

	d.Token()

	var klines []Kline
	for d.More() {
		var k Kline
		if err := d.Decode(&k); err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, errors.New(fmt.Sprint("GetKline Failed, ", err))
		}
		klines = append(klines, k)
	}

	return klines, hdr, nil
}

func (p *PriceItem) UnmarshalJSON(data []byte) error {
	type tmp struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	var raw tmp
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Symbol = Symbol(raw.Symbol)
	p.Price, _ = decimal.Parse(raw.Price)
	return nil
}

func (e Binance) FetchPrice() ([]PriceItem, http.Header, error) {
	body, hdr, err := e.request("https://api.binance.com/api/v3/ticker/price",
		map[string]string{})
	if err != nil {
		return nil, hdr, errors.New(fmt.Sprint("GetPrices Failed, ", err))
	}

	d := json.NewDecoder(body)

	var prices []PriceItem
	if err := d.Decode(&prices); err != nil {
		return nil, hdr, errors.New(fmt.Sprint("GetPrices Failed, ", err))
	}

	return prices, hdr, nil
}

func (e Binance) Ping() (http.Header, error) {
	body, hdr, err := e.request("https://api.binance.com/api/v3/ping",
		map[string]string{})
	if err != nil {
		return nil, errors.New(fmt.Sprint("Ping Failed, ", err))
	}
	body.Close()
	return hdr, nil
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type ExchangeInfo struct {
	RateLimits []RateLimit `json:"rateLimits"`
}

func (e Binance) RequestLimit() (int, http.Header, error) {
	body, hdr, err := e.request("https://api.binance.com/api/v3/exchangeInfo",
		map[string]string{"symbol": "BTCUSDT"})
	if err != nil {
		return 0, nil, errors.New(fmt.Sprint("GetMarkets Failed, ", err))
	}
	defer body.Close()

	var info ExchangeInfo
	json.NewDecoder(body).Decode(&info)

	limit := 0
	for _, item := range info.RateLimits {
		if item.RateLimitType == "REQUEST_WEIGHT" && item.Interval == "MINUTE" {
			limit = item.Limit
			break
		}
	}

	err = nil
	if limit == 0 {
		err = errors.New("could not the the request filed")
	}

	return limit, hdr, err
}

func (e Binance) request(req string, params map[string]string) (
	io.ReadCloser, http.Header, error) {

	u, err := url.Parse(req)
	if err != nil {
		return nil, nil, err
	}

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	r, err := http.Get(u.String())
	if err != nil {
		return nil, nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, nil, errors.New(fmt.Sprint("request Failed, status: , ", r.StatusCode, r.Status))
	}

	return r.Body, r.Header, nil
}
