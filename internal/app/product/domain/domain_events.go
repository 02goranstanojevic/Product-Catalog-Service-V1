package domain

import "time"

type DomainEvent interface {
	EventType() string
	AggregateID() string
}

type ProductCreatedEvent struct {
	ProductID  string
	Name       string
	Category   string
	BasePrice  string
	OccurredAt time.Time
}

func (e *ProductCreatedEvent) EventType() string   { return "product.created" }
func (e *ProductCreatedEvent) AggregateID() string { return e.ProductID }

type ProductUpdatedEvent struct {
	ProductID string
	Name      string
	Category  string
}

func (e *ProductUpdatedEvent) EventType() string   { return "product.updated" }
func (e *ProductUpdatedEvent) AggregateID() string { return e.ProductID }

type ProductActivatedEvent struct {
	ProductID string
}

func (e *ProductActivatedEvent) EventType() string   { return "product.activated" }
func (e *ProductActivatedEvent) AggregateID() string { return e.ProductID }

type ProductDeactivatedEvent struct {
	ProductID string
}

func (e *ProductDeactivatedEvent) EventType() string   { return "product.deactivated" }
func (e *ProductDeactivatedEvent) AggregateID() string { return e.ProductID }

type DiscountAppliedEvent struct {
	ProductID  string
	Percentage float64
	StartDate  time.Time
	EndDate    time.Time
}

func (e *DiscountAppliedEvent) EventType() string   { return "discount.applied" }
func (e *DiscountAppliedEvent) AggregateID() string { return e.ProductID }

type DiscountRemovedEvent struct {
	ProductID string
}

func (e *DiscountRemovedEvent) EventType() string   { return "discount.removed" }
func (e *DiscountRemovedEvent) AggregateID() string { return e.ProductID }
