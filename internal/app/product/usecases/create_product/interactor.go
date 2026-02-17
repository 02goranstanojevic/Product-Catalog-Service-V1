package create_product

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/product-catalog-service/internal/app/product/contracts"
	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/product-catalog-service/internal/pkg/clock"
	"github.com/product-catalog-service/internal/pkg/committer"
)

type Request struct {
	Name        string
	Description string
	Category    string
	Numerator   int64
	Denominator int64
}

type Interactor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer *committer.Committer
	clock     clock.Clock
}

func New(repo contracts.ProductRepository, outbox contracts.OutboxRepository, c *committer.Committer, clk clock.Clock) *Interactor {
	return &Interactor{
		repo:      repo,
		outbox:    outbox,
		committer: c,
		clock:     clk,
	}
}

func (it *Interactor) Execute(ctx context.Context, req Request) (string, error) {
	id := uuid.New().String()
	basePrice := domain.NewMoney(req.Numerator, req.Denominator)
	now := it.clock.Now()

	product, err := domain.NewProduct(id, req.Name, req.Description, req.Category, basePrice, now)
	if err != nil {
		return "", err
	}

	plan := committer.NewPlan()

	if mut := it.repo.InsertMut(product); mut != nil {
		plan.Add(mut)
	}

	for _, event := range product.DomainEvents() {
		payload, _ := json.Marshal(event)
		outboxMut := it.outbox.InsertMut(&contracts.OutboxEntry{
			EventType:   event.EventType(),
			AggregateID: event.AggregateID(),
			Payload:     payload,
		})
		plan.Add(outboxMut)
	}

	if err := it.committer.Apply(ctx, plan); err != nil {
		return "", err
	}

	return product.ID(), nil
}
