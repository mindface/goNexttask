package domain

import "errors"

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrProductionOrderNotFound = errors.New("production order not found")
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrInvalidSchedule = errors.New("invalid schedule")
)