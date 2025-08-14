package domain

import "context"

type ProductionOrderRepository interface {
	Save(ctx context.Context, order *ProductionOrder) error
	FindByID(ctx context.Context, id ProductionOrderID) (*ProductionOrder, error)
	FindAll(ctx context.Context) ([]*ProductionOrder, error)
	Update(ctx context.Context, order *ProductionOrder) error
	Delete(ctx context.Context, id ProductionOrderID) error
}