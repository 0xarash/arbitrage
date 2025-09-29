package arbitrage

import (
	"github.com/0xarash/arbitrage/src/decimal"
)

type Currency string

var USDT = Currency("USDT")

type Money struct {
	Currency Currency
	Amount   decimal.ExtendedDecimal
}

type Pair struct {
	Base  Currency
	Quote Currency
}

func (p Pair) Symbol() Symbol {
	return Symbol(p.Base + p.Quote)
}

type Symbol string

type Pairs map[Symbol]Pair
