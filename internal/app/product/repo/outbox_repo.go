package repo

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"github.com/product-catalog-service/internal/app/product/contracts"
	"github.com/product-catalog-service/internal/models/m_outbox"
	"github.com/product-catalog-service/internal/pkg/clock"
)

type OutboxRepo struct {
	clock clock.Clock
}

func NewOutboxRepo(clk clock.Clock) *OutboxRepo {
	return &OutboxRepo{clock: clk}
}

func (r *OutboxRepo) InsertMut(entry *contracts.OutboxEntry) *spanner.Mutation {
	payload := string(entry.Payload)
	if !json.Valid(entry.Payload) {
		payload = "{}"
	}

	return spanner.Insert(m_outbox.TableName, m_outbox.AllColumns, []interface{}{
		uuid.New().String(),
		entry.EventType,
		entry.AggregateID,
		payload,
		"PENDING",
		r.clock.Now(),
		(*time.Time)(nil),
	})
}
