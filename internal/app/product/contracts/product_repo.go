package contracts

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, tx *spanner.ReadOnlyTransaction, id string) (*domain.Product, error)
	InsertMut(product *domain.Product) *spanner.Mutation
	UpdateMut(product *domain.Product) *spanner.Mutation
}

type OutboxRepository interface {
	InsertMut(event *OutboxEntry) *spanner.Mutation
}

type OutboxEntry struct {
	EventType   string
	AggregateID string
	Payload     []byte
}
