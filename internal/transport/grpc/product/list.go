package product

import (
	"context"

	"github.com/product-catalog-service/internal/app/product/queries/list_products"
	pb "github.com/product-catalog-service/proto/product/v1"
)

func (h *Handler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsReply, error) {
	result, err := h.queries.ListProducts.Execute(ctx, list_products.Request{
		Category:   req.Category,
		PageSize:   req.PageSize,
		PageOffset: req.PageOffset,
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	products := make([]*pb.ProductView, 0, len(result.Products))
	for _, p := range result.Products {
		products = append(products, mapDTOToProto(p))
	}

	return &pb.ListProductsReply{
		Products:   products,
		TotalCount: result.TotalCount,
	}, nil
}
