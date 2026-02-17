package repo

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/product-catalog-service/internal/models/m_product"
	"github.com/product-catalog-service/internal/pkg/clock"
)

type ProductRepo struct {
	clock clock.Clock
}

func NewProductRepo(clk clock.Clock) *ProductRepo {
	return &ProductRepo{clock: clk}
}

func (r *ProductRepo) GetByID(ctx context.Context, tx *spanner.ReadOnlyTransaction, id string) (*domain.Product, error) {
	row, err := tx.ReadRow(ctx, m_product.TableName, spanner.Key{id}, m_product.AllColumns)
	if err != nil {
		if spanner.ErrCode(err) == 5 { // NOT_FOUND
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("reading product: %w", err)
	}

	var data m_product.ProductData
	if err := row.ToStruct(&data); err != nil {
		return nil, fmt.Errorf("scanning product: %w", err)
	}

	return toDomain(&data), nil
}

func (r *ProductRepo) InsertMut(p *domain.Product) *spanner.Mutation {
	cols := m_product.AllColumns
	discountPercent := spanner.NullNumeric{Valid: false}
	var discountStart, discountEnd *time.Time

	if d := p.Discount(); d != nil {
		discountPercent = spanner.NullNumeric{Numeric: floatToRat(d.Percentage()), Valid: true}
		s := d.StartDate()
		discountStart = &s
		e := d.EndDate()
		discountEnd = &e
	}

	return spanner.Insert(m_product.TableName, cols, []interface{}{
		p.ID(),
		p.Name(),
		p.Description(),
		p.Category(),
		p.BasePrice().Numerator(),
		p.BasePrice().Denominator(),
		discountPercent,
		discountStart,
		discountEnd,
		string(p.Status()),
		p.CreatedAt(),
		p.UpdatedAt(),
		p.ArchivedAt(),
	})
}

func (r *ProductRepo) UpdateMut(p *domain.Product) *spanner.Mutation {
	updates := map[string]interface{}{
		m_product.ProductID: p.ID(),
	}

	if p.Changes().Dirty(domain.FieldName) {
		updates[m_product.Name] = p.Name()
	}
	if p.Changes().Dirty(domain.FieldDescription) {
		updates[m_product.Description] = p.Description()
	}
	if p.Changes().Dirty(domain.FieldCategory) {
		updates[m_product.Category] = p.Category()
	}
	if p.Changes().Dirty(domain.FieldBasePrice) {
		updates[m_product.BasePriceNumerator] = p.BasePrice().Numerator()
		updates[m_product.BasePriceDenominator] = p.BasePrice().Denominator()
	}
	if p.Changes().Dirty(domain.FieldDiscount) {
		if d := p.Discount(); d != nil {
			updates[m_product.DiscountPercent] = spanner.NullNumeric{Numeric: floatToRat(d.Percentage()), Valid: true}
			s := d.StartDate()
			updates[m_product.DiscountStartDate] = &s
			e := d.EndDate()
			updates[m_product.DiscountEndDate] = &e
		} else {
			updates[m_product.DiscountPercent] = spanner.NullNumeric{Valid: false}
			updates[m_product.DiscountStartDate] = (*time.Time)(nil)
			updates[m_product.DiscountEndDate] = (*time.Time)(nil)
		}
	}
	if p.Changes().Dirty(domain.FieldStatus) {
		updates[m_product.Status] = string(p.Status())
	}
	if p.Changes().Dirty(domain.FieldArchivedAt) {
		updates[m_product.ArchivedAt] = p.ArchivedAt()
	}

	if len(updates) <= 1 {
		return nil
	}

	updates[m_product.UpdatedAt] = r.clock.Now()

	cols := make([]string, 0, len(updates))
	vals := make([]interface{}, 0, len(updates))
	for k, v := range updates {
		cols = append(cols, k)
		vals = append(vals, v)
	}

	return spanner.Update(m_product.TableName, cols, vals)
}

func toDomain(data *m_product.ProductData) *domain.Product {
	basePrice := domain.NewMoney(data.BasePriceNumerator, data.BasePriceDenominator)

	var discount *domain.Discount
	if data.DiscountPercent.Valid && data.DiscountStartDate != nil && data.DiscountEndDate != nil {
		pct, _ := data.DiscountPercent.Numeric.Float64()
		discount, _ = domain.NewDiscount(pct, *data.DiscountStartDate, *data.DiscountEndDate)
	}

	return domain.RehydrateProduct(
		data.ProductID,
		data.Name,
		data.Description,
		data.Category,
		basePrice,
		discount,
		domain.ProductStatus(data.Status),
		data.CreatedAt,
		data.UpdatedAt,
		data.ArchivedAt,
	)
}

func floatToRat(value float64) big.Rat {
	var r big.Rat
	if _, ok := r.SetString(strconv.FormatFloat(value, 'f', -1, 64)); ok {
		return r
	}
	r.SetFloat64(value)
	return r
}
