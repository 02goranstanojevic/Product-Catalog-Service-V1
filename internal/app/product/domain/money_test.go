package domain_test

import (
	"testing"

	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewMoney(t *testing.T) {
	m := domain.NewMoney(1999, 100)

	assert.Equal(t, int64(1999), m.Numerator())
	assert.Equal(t, int64(100), m.Denominator())
	assert.Equal(t, "19.99", m.String())
}

func TestMoney_Equals(t *testing.T) {
	m1 := domain.NewMoney(1000, 100)
	m2 := domain.NewMoney(1000, 100)
	m3 := domain.NewMoney(2000, 100)

	assert.True(t, m1.Equals(m2))
	assert.False(t, m1.Equals(m3))
	assert.False(t, m1.Equals(nil))
}

func TestMoney_ZeroDenominator(t *testing.T) {
	m := domain.NewMoney(100, 0)
	assert.Equal(t, int64(1), m.Denominator())
}
