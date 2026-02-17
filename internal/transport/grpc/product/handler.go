package product

import (
	"github.com/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/product-catalog-service/internal/app/product/queries/list_products"
	"github.com/product-catalog-service/internal/app/product/usecases/activate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/apply_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/archive_product"
	"github.com/product-catalog-service/internal/app/product/usecases/create_product"
	"github.com/product-catalog-service/internal/app/product/usecases/deactivate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/remove_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/update_product"
	pb "github.com/product-catalog-service/proto/product/v1"
)

type Commands struct {
	CreateProduct     *create_product.Interactor
	UpdateProduct     *update_product.Interactor
	ActivateProduct   *activate_product.Interactor
	DeactivateProduct *deactivate_product.Interactor
	ArchiveProduct    *archive_product.Interactor
	ApplyDiscount     *apply_discount.Interactor
	RemoveDiscount    *remove_discount.Interactor
}

type Queries struct {
	GetProduct   *get_product.Query
	ListProducts *list_products.Query
}

type Handler struct {
	pb.UnimplementedProductServiceServer
	commands Commands
	queries  Queries
}

func NewHandler(commands Commands, queries Queries) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}
