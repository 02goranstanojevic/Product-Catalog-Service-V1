package product

import (
	"github.com/product-catalog-service/internal/app/product/contracts"
	pb "github.com/product-catalog-service/proto/product/v1"
)

func mapDTOToProto(dto *contracts.ProductDTO) *pb.ProductView {
	return &pb.ProductView{
		Id:              dto.ID,
		Name:            dto.Name,
		Description:     dto.Description,
		Category:        dto.Category,
		BasePrice:       dto.BasePrice,
		EffectivePrice:  dto.EffectivePrice,
		DiscountPercent: dto.DiscountPercent,
		Status:          dto.Status,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
