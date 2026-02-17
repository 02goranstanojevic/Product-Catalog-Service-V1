package m_outbox

import "time"

type OutboxEventData struct {
	EventID     string     `spanner:"event_id"`
	EventType   string     `spanner:"event_type"`
	AggregateID string     `spanner:"aggregate_id"`
	Payload     string     `spanner:"payload"`
	Status      string     `spanner:"status"`
	CreatedAt   time.Time  `spanner:"created_at"`
	ProcessedAt *time.Time `spanner:"processed_at"`
}
