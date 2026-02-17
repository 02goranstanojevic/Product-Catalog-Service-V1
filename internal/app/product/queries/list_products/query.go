package list_products

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/contracts"
)

type Request struct {
	Category   string
	PageSize   int32
	PageOffset int32
}

type Response struct {
	Products   []*contracts.ProductDTO
	TotalCount int32
}

type Query struct {
	readModel contracts.ProductReadModel
	client    *spanner.Client
}

func New(readModel contracts.ProductReadModel, client *spanner.Client) *Query {
	return &Query{readModel: readModel, client: client}
}

func (q *Query) Execute(ctx context.Context, req Request) (*Response, error) {
	tx := q.client.ReadOnlyTransaction()
	defer tx.Close()

	products, total, err := q.readModel.List(ctx, tx, contracts.ListProductsFilter{
		Category:   req.Category,
		PageSize:   req.PageSize,
		PageOffset: req.PageOffset,
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		Products:   products,
		TotalCount: total,
	}, nil
}
