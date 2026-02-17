package product

import (
	"context"

	"github.com/product-catalog-service/internal/app/product/usecases/activate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/apply_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/deactivate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/remove_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/update_product"
	pb "github.com/product-catalog-service/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	err := h.commands.UpdateProduct.Execute(ctx, update_product.Request{
		ProductID:   req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.UpdateProductReply{}, nil
}

func (h *Handler) ActivateProduct(ctx context.Context, req *pb.ActivateProductRequest) (*pb.ActivateProductReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if err := h.commands.ActivateProduct.Execute(ctx, activate_product.Request{ProductID: req.ProductId}); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.ActivateProductReply{}, nil
}

func (h *Handler) DeactivateProduct(ctx context.Context, req *pb.DeactivateProductRequest) (*pb.DeactivateProductReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if err := h.commands.DeactivateProduct.Execute(ctx, deactivate_product.Request{ProductID: req.ProductId}); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.DeactivateProductReply{}, nil
}

func (h *Handler) ApplyDiscount(ctx context.Context, req *pb.ApplyDiscountRequest) (*pb.ApplyDiscountReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}
	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start_date and end_date are required")
	}

	err := h.commands.ApplyDiscount.Execute(ctx, apply_discount.Request{
		ProductID:  req.ProductId,
		Percentage: req.Percentage,
		StartDate:  req.StartDate.AsTime(),
		EndDate:    req.EndDate.AsTime(),
	})
	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.ApplyDiscountReply{}, nil
}

func (h *Handler) RemoveDiscount(ctx context.Context, req *pb.RemoveDiscountRequest) (*pb.RemoveDiscountReply, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if err := h.commands.RemoveDiscount.Execute(ctx, remove_discount.Request{ProductID: req.ProductId}); err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &pb.RemoveDiscountReply{}, nil
}
