package product

import (
	"errors"
	"fmt"

	"github.com/product-catalog-service/internal/app/product/domain"
	pb "github.com/product-catalog-service/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateCreateRequest(req *pb.CreateProductRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Category == "" {
		return fmt.Errorf("category is required")
	}
	if req.BasePriceNumerator <= 0 {
		return fmt.Errorf("base_price_numerator must be positive")
	}
	if req.BasePriceDenominator <= 0 {
		return fmt.Errorf("base_price_denominator must be positive")
	}
	return nil
}

func mapDomainErrorToGRPC(err error) error {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrEmptyProductName),
		errors.Is(err, domain.ErrEmptyCategory),
		errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrInvalidDiscountPercent),
		errors.Is(err, domain.ErrInvalidDiscountPeriod):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrProductNotActive),
		errors.Is(err, domain.ErrProductAlreadyActive),
		errors.Is(err, domain.ErrProductArchived),
		errors.Is(err, domain.ErrActiveDiscountExists),
		errors.Is(err, domain.ErrNoDiscountToRemove):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
