package get_product

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/contracts"
)

type Query struct {
	readModel contracts.ProductReadModel
	client    *spanner.Client
}

func New(readModel contracts.ProductReadModel, client *spanner.Client) *Query {
	return &Query{readModel: readModel, client: client}
}

func (q *Query) Execute(ctx context.Context, productID string) (*contracts.ProductDTO, error) {
	tx := q.client.Single()
	defer tx.Close()

	return q.readModel.GetByID(ctx, tx, productID)
}
