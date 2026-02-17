package repo

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/contracts"
	"github.com/product-catalog-service/internal/models/m_product"
	"github.com/product-catalog-service/internal/pkg/clock"
)

type ProductReadModelImpl struct {
	clock clock.Clock
}

func NewProductReadModel(clk clock.Clock) *ProductReadModelImpl {
	return &ProductReadModelImpl{clock: clk}
}

func (rm *ProductReadModelImpl) GetByID(ctx context.Context, tx *spanner.ReadOnlyTransaction, id string) (*contracts.ProductDTO, error) {
	row, err := tx.ReadRow(ctx, m_product.TableName, spanner.Key{id}, m_product.AllColumns)
	if err != nil {
		return nil, fmt.Errorf("reading product: %w", err)
	}

	var data m_product.ProductData
	if err := row.ToStruct(&data); err != nil {
		return nil, fmt.Errorf("scanning product: %w", err)
	}

	return toDTO(&data, rm.clock.Now()), nil
}

func (rm *ProductReadModelImpl) List(ctx context.Context, tx *spanner.ReadOnlyTransaction, filter contracts.ListProductsFilter) ([]*contracts.ProductDTO, int32, error) {
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := filter.PageOffset
	if offset < 0 {
		offset = 0
	}

	stmt := spanner.Statement{
		SQL: `SELECT product_id, name, description, category,
				base_price_numerator, base_price_denominator,
				discount_percent, discount_start_date, discount_end_date,
				status, created_at, updated_at, archived_at
			  FROM products
			  WHERE status = 'ACTIVE'`,
		Params: map[string]interface{}{},
	}

	if filter.Category != "" {
		stmt.SQL += ` AND category = @category`
		stmt.Params["category"] = filter.Category
	}

	countStmt := spanner.Statement{
		SQL:    fmt.Sprintf("SELECT COUNT(*) FROM (%s)", stmt.SQL),
		Params: stmt.Params,
	}

	var totalCount int32
	countIter := tx.Query(ctx, countStmt)
	countRow, err := countIter.Next()
	if err == nil {
		var c int64
		if err := countRow.Columns(&c); err == nil {
			totalCount = int32(c)
		}
	}
	countIter.Stop()

	stmt.SQL += ` ORDER BY created_at DESC LIMIT @limit OFFSET @offset`
	stmt.Params["limit"] = int64(pageSize)
	stmt.Params["offset"] = int64(offset)

	iter := tx.Query(ctx, stmt)
	defer iter.Stop()

	now := rm.clock.Now()
	var results []*contracts.ProductDTO

	for {
		row, err := iter.Next()
		if err != nil {
			break
		}

		var data m_product.ProductData
		if err := row.ToStruct(&data); err != nil {
			continue
		}

		results = append(results, toDTO(&data, now))
	}

	return results, totalCount, nil
}

func toDTO(data *m_product.ProductData, now time.Time) *contracts.ProductDTO {
	basePrice := big.NewRat(data.BasePriceNumerator, data.BasePriceDenominator)
	effectivePrice := new(big.Rat).Set(basePrice)

	if data.DiscountPercent.Valid && data.DiscountStartDate != nil && data.DiscountEndDate != nil {
		if !now.Before(*data.DiscountStartDate) && now.Before(*data.DiscountEndDate) {
			fraction := new(big.Rat).Quo(
				new(big.Rat).Set(&data.DiscountPercent.Numeric),
				big.NewRat(100, 1),
			)
			discount := new(big.Rat).Mul(basePrice, fraction)
			effectivePrice = new(big.Rat).Sub(basePrice, discount)
		}
	}

	basePriceF, _ := basePrice.Float64()
	effectivePriceF, _ := effectivePrice.Float64()

	dto := &contracts.ProductDTO{
		ID:             data.ProductID,
		Name:           data.Name,
		Description:    data.Description,
		Category:       data.Category,
		BasePrice:      fmt.Sprintf("%.2f", basePriceF),
		EffectivePrice: fmt.Sprintf("%.2f", effectivePriceF),
		Status:         data.Status,
		CreatedAt:      data.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      data.UpdatedAt.Format(time.RFC3339),
	}

	if data.DiscountPercent.Valid {
		pct, _ := data.DiscountPercent.Numeric.Float64()
		dto.DiscountPercent = &pct
	}

	return dto
}
