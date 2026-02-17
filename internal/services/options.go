package services

import (
	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/product-catalog-service/internal/app/product/queries/list_products"
	"github.com/product-catalog-service/internal/app/product/repo"
	"github.com/product-catalog-service/internal/app/product/usecases/activate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/apply_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/create_product"
	"github.com/product-catalog-service/internal/app/product/usecases/deactivate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/remove_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/update_product"
	"github.com/product-catalog-service/internal/pkg/clock"
	"github.com/product-catalog-service/internal/pkg/committer"
	grpcProduct "github.com/product-catalog-service/internal/transport/grpc/product"
)

type Container struct {
	ProductHandler *grpcProduct.Handler
}

func NewContainer(spannerClient *spanner.Client) *Container {
	clk := clock.New()
	comm := committer.New(spannerClient)

	productRepo := repo.NewProductRepo(clk)
	outboxRepo := repo.NewOutboxRepo(clk)
	readModel := repo.NewProductReadModel(clk)

	commands := grpcProduct.Commands{
		CreateProduct:     create_product.New(productRepo, outboxRepo, comm, clk),
		UpdateProduct:     update_product.New(productRepo, outboxRepo, comm, spannerClient),
		ActivateProduct:   activate_product.New(productRepo, outboxRepo, comm, spannerClient),
		DeactivateProduct: deactivate_product.New(productRepo, outboxRepo, comm, spannerClient),
		ApplyDiscount:     apply_discount.New(productRepo, outboxRepo, comm, spannerClient, clk),
		RemoveDiscount:    remove_discount.New(productRepo, outboxRepo, comm, spannerClient),
	}

	queries := grpcProduct.Queries{
		GetProduct:   get_product.New(readModel, spannerClient),
		ListProducts: list_products.New(readModel, spannerClient),
	}

	return &Container{
		ProductHandler: grpcProduct.NewHandler(commands, queries),
	}
}
