package product

import (
	"context"

	"github.com/product-catalog-service/internal/app/product/usecases/create_product"
	pb "github.com/product-catalog-service/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductReply, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	productID, err := h.commands.CreateProduct.Execute(ctx, create_product.Request{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Numerator:   req.BasePriceNumerator,
		Denominator: req.BasePriceDenominator,
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.CreateProductReply{ProductId: productID}, nil
}
