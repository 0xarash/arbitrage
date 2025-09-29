package decimal

import (
	"github.com/govalues/decimal"
)

var (
	Zero = MustNew(0, 0)
	One  = MustNew(1, 0)
)

type Kind int

const (
	Finite Kind = iota
	Infinite
)

type ExtendedDecimal struct {
	value decimal.Decimal
	kind  Kind
}

func New(value int64, scale int) (ExtendedDecimal, error) {
	v, err := decimal.New(value, scale)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func MustNew(value int64, scale int) ExtendedDecimal {
	return ExtendedDecimal{value: decimal.MustNew(value, scale), kind: Finite}
}

func NewFromFloat64(f float64, kind Kind) (ExtendedDecimal, error) {
	v, err := decimal.NewFromFloat64(f)
	return ExtendedDecimal{value: v, kind: kind}, err
}

func Parse(s string) (ExtendedDecimal, error) {
	v, err := decimal.Parse(s)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func Mean(d ...ExtendedDecimal) (ExtendedDecimal, error) {
	decs := make([]decimal.Decimal, len(d))
	for i, v := range d {
		decs[i] = v.value
	}
	v, err := decimal.Mean(decs...)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Mul(e ExtendedDecimal) (ExtendedDecimal, error) {
	v, err := d.value.Mul(e.value)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Quo(e ExtendedDecimal) (ExtendedDecimal, error) {
	v, err := d.value.Quo(e.value)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Add(e ExtendedDecimal) (ExtendedDecimal, error) {
	if d.kind == Infinite || e.kind == Infinite {
		return ExtendedDecimal{value: decimal.Decimal{}, kind: Infinite}, nil
	}
	v, err := d.value.Add(e.value)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Sub(e ExtendedDecimal) (ExtendedDecimal, error) {
	if d.kind == Infinite || e.kind == Infinite {
		return ExtendedDecimal{value: decimal.Decimal{}, kind: Infinite}, nil
	}
	v, err := d.value.Sub(e.value)
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Neg() ExtendedDecimal {
	return ExtendedDecimal{value: d.value.Neg(), kind: Finite}
}

func (d ExtendedDecimal) Log() (ExtendedDecimal, error) {
	v, err := d.value.Log()
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Inv() (ExtendedDecimal, error) {
	v, err := d.value.Inv()
	return ExtendedDecimal{value: v, kind: Finite}, err
}

func (d ExtendedDecimal) Less(e ExtendedDecimal) bool {
	switch d.kind {
	case Infinite:
		return false
	case Finite:
		switch e.kind {
		case Infinite:
			return true
		case Finite:
			return d.value.Less(e.value)
		}
	}
	return false
}

func (d ExtendedDecimal) IsZero() bool {
	return d.value.IsZero()
}

func (d ExtendedDecimal) CmpTotal(e ExtendedDecimal) int {
	return d.value.CmpTotal(e.value)
}

func (d ExtendedDecimal) String() string {
	return d.value.String()
}
