package domain

import "errors"

var (
	ErrEmptyProductName       = errors.New("product name cannot be empty")
	ErrEmptyCategory          = errors.New("product category cannot be empty")
	ErrInvalidPrice           = errors.New("base price must be positive")
	ErrProductNotActive       = errors.New("product is not active")
	ErrProductAlreadyActive   = errors.New("product is already active")
	ErrProductArchived        = errors.New("product is archived")
	ErrProductNotFound        = errors.New("product not found")
	ErrInvalidDiscountPercent = errors.New("discount percentage must be between 0 and 100")
	ErrInvalidDiscountPeriod  = errors.New("discount period is invalid")
	ErrActiveDiscountExists   = errors.New("product already has an active discount")
	ErrNoDiscountToRemove     = errors.New("product has no discount to remove")
)
