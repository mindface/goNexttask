package domain

import "time"

type EventType string

const (
	EventProductionOrderCreated   EventType = "ProductionOrderCreated"
	EventProductionOrderStarted   EventType = "ProductionOrderStarted"
	EventProductionOrderCompleted EventType = "ProductionOrderCompleted"
	EventProductionOrderDelayed   EventType = "ProductionOrderDelayed"
	EventProductionOrderCancelled EventType = "ProductionOrderCancelled"
)

type DomainEvent interface {
	GetEventType() EventType
	GetOccurredAt() time.Time
	GetAggregateID() string
}

type ProductionOrderEvent struct {
	EventType   EventType
	OrderID     ProductionOrderID
	OrderNumber string
	OccurredAt  time.Time
	Payload     map[string]interface{}
}

func (e ProductionOrderEvent) GetEventType() EventType {
	return e.EventType
}

func (e ProductionOrderEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e ProductionOrderEvent) GetAggregateID() string {
	return string(e.OrderID)
}

func NewProductionOrderCreatedEvent(order *ProductionOrder) DomainEvent {
	return ProductionOrderEvent{
		EventType:   EventProductionOrderCreated,
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		OccurredAt:  time.Now(),
		Payload: map[string]interface{}{
			"partID":   order.PartID,
			"quantity": order.Quantity,
			"status":   order.Status,
		},
	}
}

func NewProductionOrderStartedEvent(order *ProductionOrder) DomainEvent {
	return ProductionOrderEvent{
		EventType:   EventProductionOrderStarted,
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		OccurredAt:  time.Now(),
		Payload: map[string]interface{}{
			"status": order.Status,
		},
	}
}