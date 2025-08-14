package domain

import (
	"context"
	"time"
)

type ProductionSchedulingService struct {
	repo ProductionOrderRepository
}

func NewProductionSchedulingService(repo ProductionOrderRepository) *ProductionSchedulingService {
	return &ProductionSchedulingService{
		repo: repo,
	}
}

func (s *ProductionSchedulingService) ScheduleProduction(
	ctx context.Context,
	orderNumber string,
	partID PartID,
	quantity int,
	plannedStart time.Time,
	plannedEnd time.Time,
	machineIDs []MachineID,
) (*ProductionOrder, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	
	if plannedEnd.Before(plannedStart) {
		return nil, ErrInvalidSchedule
	}
	
	schedule := Schedule{
		PlannedStart:     plannedStart,
		PlannedEnd:       plannedEnd,
		AssignedMachines: machineIDs,
	}
	
	order := NewProductionOrder(orderNumber, partID, quantity, schedule)
	
	if err := s.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	
	return order, nil
}

func (s *ProductionSchedulingService) OptimizeSchedule(ctx context.Context, orders []*ProductionOrder) error {
	// TODO: 実装予定 - リソース可用性、NC機の稼働状況を考慮した最適化ロジック
	return nil
}