package product

import (
	"context"

	pb "github.com/product-catalog-service/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	dto, err := h.queries.GetProduct.Execute(ctx, req.ProductId)
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.GetProductReply{Product: mapDTOToProto(dto)}, nil
}
