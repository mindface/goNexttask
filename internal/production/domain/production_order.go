package domain

import (
	"time"
)

type ProductionOrderStatus string

const (
	StatusPlanned    ProductionOrderStatus = "planned"
	StatusInProgress ProductionOrderStatus = "in_progress"
	StatusCompleted  ProductionOrderStatus = "completed"
	StatusDelayed    ProductionOrderStatus = "delayed"
	StatusCancelled  ProductionOrderStatus = "cancelled"
)

type ProductionOrderID string

type PartID string

type MachineID string

type ProductionOrder struct {
	ID          ProductionOrderID
	OrderNumber string
	PartID      PartID
	Quantity    int
	Status      ProductionOrderStatus
	Schedule    Schedule
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Schedule struct {
	PlannedStart     time.Time
	PlannedEnd       time.Time
	AssignedMachines []MachineID
}

func NewProductionOrder(orderNumber string, partID PartID, quantity int, schedule Schedule) *ProductionOrder {
	now := time.Now()
	return &ProductionOrder{
		ID:          ProductionOrderID("order-" + orderNumber),
		OrderNumber: orderNumber,
		PartID:      partID,
		Quantity:    quantity,
		Status:      StatusPlanned,
		Schedule:    schedule,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (po *ProductionOrder) Start() error {
	if po.Status != StatusPlanned {
		return ErrInvalidStateTransition
	}
	po.Status = StatusInProgress
	po.UpdatedAt = time.Now()
	return nil
}

func (po *ProductionOrder) Complete() error {
	if po.Status != StatusInProgress {
		return ErrInvalidStateTransition
	}
	po.Status = StatusCompleted
	po.UpdatedAt = time.Now()
	return nil
}

func (po *ProductionOrder) Delay() error {
	if po.Status != StatusInProgress && po.Status != StatusPlanned {
		return ErrInvalidStateTransition
	}
	po.Status = StatusDelayed
	po.UpdatedAt = time.Now()
	return nil
}

func (po *ProductionOrder) Cancel() error {
	if po.Status == StatusCompleted || po.Status == StatusCancelled {
		return ErrInvalidStateTransition
	}
	po.Status = StatusCancelled
	po.UpdatedAt = time.Now()
	return nil
}