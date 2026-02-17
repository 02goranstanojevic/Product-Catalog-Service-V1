package domain

import "time"

type Discount struct {
	percentage float64
	startDate  time.Time
	endDate    time.Time
}

func NewDiscount(percentage float64, startDate, endDate time.Time) (*Discount, error) {
	if percentage <= 0 || percentage > 100 {
		return nil, ErrInvalidDiscountPercent
	}
	if !endDate.After(startDate) {
		return nil, ErrInvalidDiscountPeriod
	}

	return &Discount{
		percentage: percentage,
		startDate:  startDate,
		endDate:    endDate,
	}, nil
}

func (d *Discount) IsValidAt(t time.Time) bool {
	return !t.Before(d.startDate) && t.Before(d.endDate)
}

func (d *Discount) Percentage() float64  { return d.percentage }
func (d *Discount) StartDate() time.Time { return d.startDate }
func (d *Discount) EndDate() time.Time   { return d.endDate }
