package domain

import (
	"math/big"
	"time"
)

type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "DRAFT"
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
	ProductStatusArchived ProductStatus = "ARCHIVED"
)

const (
	FieldName        = "name"
	FieldDescription = "description"
	FieldCategory    = "category"
	FieldBasePrice   = "base_price"
	FieldDiscount    = "discount"
	FieldStatus      = "status"
	FieldArchivedAt  = "archived_at"
)

type Product struct {
	id          string
	name        string
	description string
	category    string
	basePrice   *Money
	discount    *Discount
	status      ProductStatus
	createdAt   time.Time
	updatedAt   time.Time
	archivedAt  *time.Time
	changes     *ChangeTracker
	events      []DomainEvent
}

func NewProduct(id, name, description, category string, basePrice *Money, now time.Time) (*Product, error) {
	if name == "" {
		return nil, ErrEmptyProductName
	}
	if category == "" {
		return nil, ErrEmptyCategory
	}
	if basePrice == nil || basePrice.Amount().Sign() <= 0 {
		return nil, ErrInvalidPrice
	}

	p := &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		status:      ProductStatusDraft,
		createdAt:   now,
		updatedAt:   now,
		changes:     NewChangeTracker(),
		events:      make([]DomainEvent, 0),
	}

	p.events = append(p.events, &ProductCreatedEvent{
		ProductID:  id,
		Name:       name,
		Category:   category,
		BasePrice:  basePrice.String(),
		OccurredAt: now,
	})

	return p, nil
}

func RehydrateProduct(
	id, name, description, category string,
	basePrice *Money,
	discount *Discount,
	status ProductStatus,
	createdAt, updatedAt time.Time,
	archivedAt *time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		category:    category,
		basePrice:   basePrice,
		discount:    discount,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		archivedAt:  archivedAt,
		changes:     NewChangeTracker(),
		events:      make([]DomainEvent, 0),
	}
}

func (p *Product) UpdateDetails(name, description, category string) error {
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}

	if name != "" && name != p.name {
		p.name = name
		p.changes.MarkDirty(FieldName)
	}
	if description != p.description {
		p.description = description
		p.changes.MarkDirty(FieldDescription)
	}
	if category != "" && category != p.category {
		p.category = category
		p.changes.MarkDirty(FieldCategory)
	}

	if p.changes.HasChanges() {
		p.events = append(p.events, &ProductUpdatedEvent{
			ProductID: p.id,
			Name:      p.name,
			Category:  p.category,
		})
	}

	return nil
}

func (p *Product) Activate() error {
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}
	if p.status == ProductStatusActive {
		return ErrProductAlreadyActive
	}

	p.status = ProductStatusActive
	p.changes.MarkDirty(FieldStatus)
	p.events = append(p.events, &ProductActivatedEvent{ProductID: p.id})
	return nil
}

func (p *Product) Deactivate() error {
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}

	p.status = ProductStatusInactive
	p.changes.MarkDirty(FieldStatus)
	p.events = append(p.events, &ProductDeactivatedEvent{ProductID: p.id})
	return nil
}

func (p *Product) Archive(now time.Time) error {
	if p.status == ProductStatusArchived {
		return ErrProductArchived
	}

	p.status = ProductStatusArchived
	p.archivedAt = &now
	p.changes.MarkDirty(FieldStatus)
	p.changes.MarkDirty(FieldArchivedAt)
	return nil
}

func (p *Product) ApplyDiscount(discount *Discount, now time.Time) error {
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}
	if !discount.IsValidAt(now) {
		return ErrInvalidDiscountPeriod
	}
	if p.discount != nil && p.discount.IsValidAt(now) {
		return ErrActiveDiscountExists
	}

	p.discount = discount
	p.changes.MarkDirty(FieldDiscount)
	p.events = append(p.events, &DiscountAppliedEvent{
		ProductID:  p.id,
		Percentage: discount.Percentage(),
		StartDate:  discount.StartDate(),
		EndDate:    discount.EndDate(),
	})
	return nil
}

func (p *Product) RemoveDiscount() error {
	if p.status != ProductStatusActive {
		return ErrProductNotActive
	}
	if p.discount == nil {
		return ErrNoDiscountToRemove
	}

	p.discount = nil
	p.changes.MarkDirty(FieldDiscount)
	p.events = append(p.events, &DiscountRemovedEvent{ProductID: p.id})
	return nil
}

func (p *Product) EffectivePrice(now time.Time) *big.Rat {
	if p.discount != nil && p.discount.IsValidAt(now) {
		discountFraction := new(big.Rat).Quo(
			new(big.Rat).SetFloat64(p.discount.Percentage()),
			big.NewRat(100, 1),
		)
		discountAmount := new(big.Rat).Mul(p.basePrice.Amount(), discountFraction)
		return new(big.Rat).Sub(p.basePrice.Amount(), discountAmount)
	}
	return new(big.Rat).Set(p.basePrice.Amount())
}

func (p *Product) ID() string                  { return p.id }
func (p *Product) Name() string                { return p.name }
func (p *Product) Description() string         { return p.description }
func (p *Product) Category() string            { return p.category }
func (p *Product) BasePrice() *Money           { return p.basePrice }
func (p *Product) Discount() *Discount         { return p.discount }
func (p *Product) Status() ProductStatus       { return p.status }
func (p *Product) CreatedAt() time.Time        { return p.createdAt }
func (p *Product) UpdatedAt() time.Time        { return p.updatedAt }
func (p *Product) ArchivedAt() *time.Time      { return p.archivedAt }
func (p *Product) Changes() *ChangeTracker     { return p.changes }
func (p *Product) DomainEvents() []DomainEvent { return p.events }
func (p *Product) ClearEvents()                { p.events = nil }
