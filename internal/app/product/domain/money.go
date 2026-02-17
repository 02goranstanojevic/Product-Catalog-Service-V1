package domain

import (
	"fmt"
	"math/big"
)

type Money struct {
	amount *big.Rat
}

func NewMoney(numerator, denominator int64) *Money {
	if denominator == 0 {
		denominator = 1
	}
	return &Money{amount: big.NewRat(numerator, denominator)}
}

func NewMoneyFromRat(r *big.Rat) *Money {
	return &Money{amount: new(big.Rat).Set(r)}
}

func (m *Money) Amount() *big.Rat {
	return new(big.Rat).Set(m.amount)
}

func (m *Money) Numerator() int64 {
	return m.amount.Num().Int64()
}

func (m *Money) Denominator() int64 {
	return m.amount.Denom().Int64()
}

func (m *Money) String() string {
	f, _ := m.amount.Float64()
	return fmt.Sprintf("%.2f", f)
}

func (m *Money) Equals(other *Money) bool {
	if other == nil {
		return false
	}
	return m.amount.Cmp(other.amount) == 0
}
