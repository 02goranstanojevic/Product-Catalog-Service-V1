package m_product

import (
	"time"

	"cloud.google.com/go/spanner"
)

type ProductData struct {
	ProductID            string              `spanner:"product_id"`
	Name                 string              `spanner:"name"`
	Description          string              `spanner:"description"`
	Category             string              `spanner:"category"`
	BasePriceNumerator   int64               `spanner:"base_price_numerator"`
	BasePriceDenominator int64               `spanner:"base_price_denominator"`
	DiscountPercent      spanner.NullNumeric `spanner:"discount_percent"`
	DiscountStartDate    *time.Time          `spanner:"discount_start_date"`
	DiscountEndDate      *time.Time          `spanner:"discount_end_date"`
	Status               string              `spanner:"status"`
	CreatedAt            time.Time           `spanner:"created_at"`
	UpdatedAt            time.Time           `spanner:"updated_at"`
	ArchivedAt           *time.Time          `spanner:"archived_at"`
}
