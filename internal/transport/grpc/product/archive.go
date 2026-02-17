package product

import (
	"context"

	"github.com/product-catalog-service/internal/app/product/usecases/archive_product"
	pb "github.com/product-catalog-service/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ArchiveProduct(ctx context.Context, req *pb.ArchiveProductRequest) (*pb.ArchiveProductReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if err := h.commands.ArchiveProduct.Execute(ctx, archive_product.Request{ProductID: req.ProductId}); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.ArchiveProductReply{}, nil
}
