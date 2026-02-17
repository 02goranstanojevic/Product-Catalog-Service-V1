package domain_test

import (
	"testing"
	"time"

	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiscount_Success(t *testing.T) {
	start := time.Now()
	end := start.Add(24 * time.Hour)

	d, err := domain.NewDiscount(20, start, end)

	require.NoError(t, err)
	assert.Equal(t, 20.0, d.Percentage())
	assert.Equal(t, start, d.StartDate())
	assert.Equal(t, end, d.EndDate())
}

func TestNewDiscount_InvalidPercentage(t *testing.T) {
	start := time.Now()
	end := start.Add(24 * time.Hour)

	_, err := domain.NewDiscount(0, start, end)
	assert.ErrorIs(t, err, domain.ErrInvalidDiscountPercent)

	_, err = domain.NewDiscount(101, start, end)
	assert.ErrorIs(t, err, domain.ErrInvalidDiscountPercent)

	_, err = domain.NewDiscount(-5, start, end)
	assert.ErrorIs(t, err, domain.ErrInvalidDiscountPercent)
}

func TestNewDiscount_InvalidPeriod(t *testing.T) {
	now := time.Now()

	_, err := domain.NewDiscount(10, now, now)
	assert.ErrorIs(t, err, domain.ErrInvalidDiscountPeriod)

	_, err = domain.NewDiscount(10, now, now.Add(-time.Hour))
	assert.ErrorIs(t, err, domain.ErrInvalidDiscountPeriod)
}

func TestDiscount_IsValidAt(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	d, _ := domain.NewDiscount(10, start, end)

	assert.True(t, d.IsValidAt(time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)))
	assert.True(t, d.IsValidAt(start))
	assert.False(t, d.IsValidAt(end))
	assert.False(t, d.IsValidAt(start.Add(-time.Second)))
}
