package services

import (
	"math/big"
	"time"

	"github.com/product-catalog-service/internal/app/product/domain"
)

type PricingCalculator struct{}

func NewPricingCalculator() *PricingCalculator {
	return &PricingCalculator{}
}

func (pc *PricingCalculator) CalculateEffectivePrice(product *domain.Product, now time.Time) *big.Rat {
	return product.EffectivePrice(now)
}

func (pc *PricingCalculator) CalculateDiscountAmount(basePrice *domain.Money, percentage float64) *big.Rat {
	fraction := new(big.Rat).Quo(
		new(big.Rat).SetFloat64(percentage),
		big.NewRat(100, 1),
	)
	return new(big.Rat).Mul(basePrice.Amount(), fraction)
}
