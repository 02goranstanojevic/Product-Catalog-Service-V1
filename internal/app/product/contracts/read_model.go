package contracts

import (
	"context"

	"cloud.google.com/go/spanner"
)

type ProductDTO struct {
	ID              string
	Name            string
	Description     string
	Category        string
	BasePrice       string
	EffectivePrice  string
	DiscountPercent *float64
	Status          string
	CreatedAt       string
	UpdatedAt       string
}

type ListProductsFilter struct {
	Category   string
	PageSize   int32
	PageOffset int32
}

type ProductReadModel interface {
	GetByID(ctx context.Context, tx *spanner.ReadOnlyTransaction, id string) (*ProductDTO, error)
	List(ctx context.Context, tx *spanner.ReadOnlyTransaction, filter ListProductsFilter) ([]*ProductDTO, int32, error)
}
