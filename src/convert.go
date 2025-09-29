package arbitrage

import (
	"errors"
	"fmt"

	"github.com/0xarash/arbitrage/src/decimal"

	"github.com/rs/zerolog/log"
)

type Converter struct {
	mapper map[Currency]*TargetCurrencyConverter
}

type TargetCurrencyConverter struct {
	baseRates map[Currency]decimal.ExtendedDecimal
}

func NewConverter(prices Prices, baseCurrencies []Currency) (Converter, error) {
	if len(prices) == 0 {
		return Converter{}, errors.New("NewConverter: Empty depth is provided")
	}

	c := Converter{mapper: make(map[Currency]*TargetCurrencyConverter)}
	for _, curr := range baseCurrencies {
		if _, ok := c.mapper[curr]; !ok {
			target := &TargetCurrencyConverter{
				baseRates: make(map[Currency]decimal.ExtendedDecimal),
			}
			for _, price := range prices {
				if price.Pair.Quote != curr {
					continue
				}

				target.baseRates[price.Pair.Base] = price.Price
			}
			c.mapper[curr] = target
		}
	}

	return c, nil
}

func (c *Converter) Convert(money Money, currency Currency) (Money, error) {
	if money.Currency == currency {
		return money, nil
	}

	conv := c.mapper[currency]
	baseRate, ok := conv.baseRates[money.Currency]
	if !ok {
		log.Warn().Msgf("Currency %v is not L1 cache .", money.Currency)
		return Money{}, nil
	}

	liquidity, err := money.Amount.Mul(baseRate)
	if err != nil {
		panic(fmt.Sprintf("Decimal multiplication failed for currency %v",
			money.Currency))
	}

	return Money{Amount: liquidity, Currency: currency}, nil
}
