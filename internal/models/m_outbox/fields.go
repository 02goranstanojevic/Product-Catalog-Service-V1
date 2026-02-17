package m_outbox

const (
	TableName = "outbox_events"

	EventID     = "event_id"
	EventType   = "event_type"
	AggregateID = "aggregate_id"
	Payload     = "payload"
	Status      = "status"
	CreatedAt   = "created_at"
	ProcessedAt = "processed_at"
)

var AllColumns = []string{
	EventID, EventType, AggregateID, Payload, Status, CreatedAt, ProcessedAt,
}
