package domain_test

import (
	"testing"
	"time"

	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProduct_Success(t *testing.T) {
	now := time.Now().UTC()
	price := domain.NewMoney(1999, 100)

	p, err := domain.NewProduct("id-1", "Widget", "A widget", "electronics", price, now)

	require.NoError(t, err)
	assert.Equal(t, "id-1", p.ID())
	assert.Equal(t, "Widget", p.Name())
	assert.Equal(t, "electronics", p.Category())
	assert.Equal(t, domain.ProductStatusDraft, p.Status())
	assert.Len(t, p.DomainEvents(), 1)
	assert.Equal(t, "product.created", p.DomainEvents()[0].EventType())
}

func TestNewProduct_EmptyName(t *testing.T) {
	price := domain.NewMoney(1000, 100)
	_, err := domain.NewProduct("id-1", "", "desc", "cat", price, time.Now())
	assert.ErrorIs(t, err, domain.ErrEmptyProductName)
}

func TestNewProduct_EmptyCategory(t *testing.T) {
	price := domain.NewMoney(1000, 100)
	_, err := domain.NewProduct("id-1", "Name", "desc", "", price, time.Now())
	assert.ErrorIs(t, err, domain.ErrEmptyCategory)
}

func TestNewProduct_InvalidPrice(t *testing.T) {
	price := domain.NewMoney(-100, 100)
	_, err := domain.NewProduct("id-1", "Name", "desc", "cat", price, time.Now())
	assert.ErrorIs(t, err, domain.ErrInvalidPrice)
}

func TestProduct_Activate(t *testing.T) {
	p := createTestProduct(t)

	err := p.Activate()

	require.NoError(t, err)
	assert.Equal(t, domain.ProductStatusActive, p.Status())
	assert.True(t, p.Changes().Dirty(domain.FieldStatus))
}

func TestProduct_Activate_AlreadyActive(t *testing.T) {
	p := createActiveProduct(t)

	err := p.Activate()
	assert.ErrorIs(t, err, domain.ErrProductAlreadyActive)
}

func TestProduct_Deactivate(t *testing.T) {
	p := createActiveProduct(t)

	err := p.Deactivate()

	require.NoError(t, err)
	assert.Equal(t, domain.ProductStatusInactive, p.Status())
}

func TestProduct_Deactivate_NotActive(t *testing.T) {
	p := createTestProduct(t)

	err := p.Deactivate()
	assert.ErrorIs(t, err, domain.ErrProductNotActive)
}

func TestProduct_Archive(t *testing.T) {
	p := createTestProduct(t)
	now := time.Now()

	err := p.Archive(now)

	require.NoError(t, err)
	assert.Equal(t, domain.ProductStatusArchived, p.Status())
	assert.NotNil(t, p.ArchivedAt())
}

func TestProduct_Archive_AlreadyArchived(t *testing.T) {
	p := createTestProduct(t)
	_ = p.Archive(time.Now())

	err := p.Archive(time.Now())
	assert.ErrorIs(t, err, domain.ErrProductArchived)
}

func TestProduct_UpdateDetails(t *testing.T) {
	p := createTestProduct(t)

	err := p.UpdateDetails("New Name", "New Desc", "new-cat")

	require.NoError(t, err)
	assert.Equal(t, "New Name", p.Name())
	assert.Equal(t, "New Desc", p.Description())
	assert.Equal(t, "new-cat", p.Category())
	assert.True(t, p.Changes().HasChanges())
}

func TestProduct_UpdateDetails_Archived(t *testing.T) {
	p := createTestProduct(t)
	_ = p.Archive(time.Now())

	err := p.UpdateDetails("name", "desc", "cat")
	assert.ErrorIs(t, err, domain.ErrProductArchived)
}

func TestProduct_ApplyDiscount(t *testing.T) {
	p := createActiveProduct(t)
	now := time.Now()
	discount, _ := domain.NewDiscount(20, now.Add(-time.Hour), now.Add(24*time.Hour))

	err := p.ApplyDiscount(discount, now)

	require.NoError(t, err)
	assert.NotNil(t, p.Discount())
	assert.True(t, p.Changes().Dirty(domain.FieldDiscount))
}

func TestProduct_ApplyDiscount_NotActive(t *testing.T) {
	p := createTestProduct(t)
	now := time.Now()
	discount, _ := domain.NewDiscount(20, now.Add(-time.Hour), now.Add(24*time.Hour))

	err := p.ApplyDiscount(discount, now)
	assert.ErrorIs(t, err, domain.ErrProductNotActive)
}

func TestProduct_ApplyDiscount_ActiveDiscountExists(t *testing.T) {
	p := createActiveProduct(t)
	now := time.Now()
	d1, _ := domain.NewDiscount(10, now.Add(-time.Hour), now.Add(24*time.Hour))
	_ = p.ApplyDiscount(d1, now)

	d2, _ := domain.NewDiscount(15, now.Add(-time.Hour), now.Add(48*time.Hour))
	err := p.ApplyDiscount(d2, now)
	assert.ErrorIs(t, err, domain.ErrActiveDiscountExists)
}

func TestProduct_RemoveDiscount(t *testing.T) {
	p := createActiveProduct(t)
	now := time.Now()
	discount, _ := domain.NewDiscount(20, now.Add(-time.Hour), now.Add(24*time.Hour))
	_ = p.ApplyDiscount(discount, now)

	err := p.RemoveDiscount()

	require.NoError(t, err)
	assert.Nil(t, p.Discount())
}

func TestProduct_RemoveDiscount_NoDiscount(t *testing.T) {
	p := createActiveProduct(t)

	err := p.RemoveDiscount()
	assert.ErrorIs(t, err, domain.ErrNoDiscountToRemove)
}

func TestProduct_EffectivePrice_NoDiscount(t *testing.T) {
	p := createActiveProduct(t)
	now := time.Now()

	price := p.EffectivePrice(now)
	f, _ := price.Float64()

	assert.InDelta(t, 19.99, f, 0.001)
}

func TestProduct_EffectivePrice_WithDiscount(t *testing.T) {
	p := createActiveProduct(t)
	now := time.Now()
	discount, _ := domain.NewDiscount(20, now.Add(-time.Hour), now.Add(24*time.Hour))
	_ = p.ApplyDiscount(discount, now)

	price := p.EffectivePrice(now)
	f, _ := price.Float64()

	assert.InDelta(t, 15.992, f, 0.001)
}

func createTestProduct(t *testing.T) *domain.Product {
	t.Helper()
	price := domain.NewMoney(1999, 100)
	p, err := domain.NewProduct("test-id", "Test Product", "A test product", "electronics", price, time.Now().UTC())
	require.NoError(t, err)
	return p
}

func createActiveProduct(t *testing.T) *domain.Product {
	t.Helper()
	p := createTestProduct(t)
	require.NoError(t, p.Activate())
	return p
}
