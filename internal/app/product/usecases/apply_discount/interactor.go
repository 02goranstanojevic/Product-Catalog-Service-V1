package apply_discount

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/contracts"
	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/product-catalog-service/internal/pkg/clock"
	"github.com/product-catalog-service/internal/pkg/committer"
)

type Request struct {
	ProductID  string
	Percentage float64
	StartDate  time.Time
	EndDate    time.Time
}

type Interactor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer *committer.Committer
	client    *spanner.Client
	clock     clock.Clock
}

func New(repo contracts.ProductRepository, outbox contracts.OutboxRepository, c *committer.Committer, client *spanner.Client, clk clock.Clock) *Interactor {
	return &Interactor{
		repo:      repo,
		outbox:    outbox,
		committer: c,
		client:    client,
		clock:     clk,
	}
}

func (it *Interactor) Execute(ctx context.Context, req Request) error {
	tx := it.client.Single()
	defer tx.Close()

	product, err := it.repo.GetByID(ctx, tx, req.ProductID)
	if err != nil {
		return err
	}

	discount, err := domain.NewDiscount(req.Percentage, req.StartDate, req.EndDate)
	if err != nil {
		return err
	}

	if err := product.ApplyDiscount(discount, it.clock.Now()); err != nil {
		return err
	}

	plan := committer.NewPlan()

	if mut := it.repo.UpdateMut(product); mut != nil {
		plan.Add(mut)
	}

	for _, event := range product.DomainEvents() {
		payload, _ := json.Marshal(event)
		plan.Add(it.outbox.InsertMut(&contracts.OutboxEntry{
			EventType:   event.EventType(),
			AggregateID: event.AggregateID(),
			Payload:     payload,
		}))
	}

	return it.committer.Apply(ctx, plan)
}
