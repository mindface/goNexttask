package domain

import "context"

type InspectionRepository interface {
	Save(ctx context.Context, inspection *Inspection) error
	FindByID(ctx context.Context, id InspectionID) (*Inspection, error)
	FindByLotNumber(ctx context.Context, lotNumber string) ([]*Inspection, error)
	FindByProductionOrderID(ctx context.Context, orderID string) ([]*Inspection, error)
	Update(ctx context.Context, inspection *Inspection) error
}