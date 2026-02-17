package update_product

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/app/product/contracts"
	"github.com/product-catalog-service/internal/pkg/committer"
)

type Request struct {
	ProductID   string
	Name        string
	Description string
	Category    string
}

type Interactor struct {
	repo      contracts.ProductRepository
	outbox    contracts.OutboxRepository
	committer *committer.Committer
	client    *spanner.Client
}

func New(repo contracts.ProductRepository, outbox contracts.OutboxRepository, c *committer.Committer, client *spanner.Client) *Interactor {
	return &Interactor{
		repo:      repo,
		outbox:    outbox,
		committer: c,
		client:    client,
	}
}

func (it *Interactor) Execute(ctx context.Context, req Request) error {
	tx := it.client.Single()
	defer tx.Close()

	product, err := it.repo.GetByID(ctx, tx, req.ProductID)
	if err != nil {
		return err
	}

	if err := product.UpdateDetails(req.Name, req.Description, req.Category); err != nil {
		return err
	}

	if !product.Changes().HasChanges() {
		return nil
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
