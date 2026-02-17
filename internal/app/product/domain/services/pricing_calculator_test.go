package services_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/product-catalog-service/internal/app/product/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPricingCalculator_EffectivePrice_NoDiscount(t *testing.T) {
	calc := services.NewPricingCalculator()
	price := domain.NewMoney(2000, 100)
	p, err := domain.NewProduct("id", "Test", "desc", "cat", price, time.Now())
	require.NoError(t, err)

	result := calc.CalculateEffectivePrice(p, time.Now())
	f, _ := result.Float64()

	assert.InDelta(t, 20.0, f, 0.001)
}

func TestPricingCalculator_DiscountAmount(t *testing.T) {
	calc := services.NewPricingCalculator()
	price := domain.NewMoney(10000, 100)

	result := calc.CalculateDiscountAmount(price, 25)
	expected := big.NewRat(25, 1)

	assert.Equal(t, 0, result.Cmp(expected))
}
