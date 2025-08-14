package application

import (
	"context"
	"goNexttask/internal/production/domain"
	"time"
)

type CreateProductionOrderInput struct {
	OrderNumber      string
	PartID           string
	Quantity         int
	PlannedStartDate time.Time
	PlannedEndDate   time.Time
	MachineIDs       []string
}

type ProductionOrderOutput struct {
	ID               string
	OrderNumber      string
	PartID           string
	Quantity         int
	Status           string
	PlannedStartDate time.Time
	PlannedEndDate   time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ProductionUseCase struct {
	repo               domain.ProductionOrderRepository
	schedulingService  *domain.ProductionSchedulingService
}

func NewProductionUseCase(repo domain.ProductionOrderRepository) *ProductionUseCase {
	return &ProductionUseCase{
		repo:              repo,
		schedulingService: domain.NewProductionSchedulingService(repo),
	}
}

func (uc *ProductionUseCase) CreateProductionOrder(ctx context.Context, input CreateProductionOrderInput) (*ProductionOrderOutput, error) {
	machineIDs := make([]domain.MachineID, len(input.MachineIDs))
	for i, id := range input.MachineIDs {
		machineIDs[i] = domain.MachineID(id)
	}
	
	order, err := uc.schedulingService.ScheduleProduction(
		ctx,
		input.OrderNumber,
		domain.PartID(input.PartID),
		input.Quantity,
		input.PlannedStartDate,
		input.PlannedEndDate,
		machineIDs,
	)
	if err != nil {
		return nil, err
	}
	
	return &ProductionOrderOutput{
		ID:               string(order.ID),
		OrderNumber:      order.OrderNumber,
		PartID:           string(order.PartID),
		Quantity:         order.Quantity,
		Status:           string(order.Status),
		PlannedStartDate: order.Schedule.PlannedStart,
		PlannedEndDate:   order.Schedule.PlannedEnd,
		CreatedAt:        order.CreatedAt,
		UpdatedAt:        order.UpdatedAt,
	}, nil
}

func (uc *ProductionUseCase) GetProductionOrder(ctx context.Context, id string) (*ProductionOrderOutput, error) {
	order, err := uc.repo.FindByID(ctx, domain.ProductionOrderID(id))
	if err != nil {
		return nil, err
	}
	
	return &ProductionOrderOutput{
		ID:               string(order.ID),
		OrderNumber:      order.OrderNumber,
		PartID:           string(order.PartID),
		Quantity:         order.Quantity,
		Status:           string(order.Status),
		PlannedStartDate: order.Schedule.PlannedStart,
		PlannedEndDate:   order.Schedule.PlannedEnd,
		CreatedAt:        order.CreatedAt,
		UpdatedAt:        order.UpdatedAt,
	}, nil
}

func (uc *ProductionUseCase) GetAllProductionOrders(ctx context.Context) ([]*ProductionOrderOutput, error) {
	orders, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	
	outputs := make([]*ProductionOrderOutput, len(orders))
	for i, order := range orders {
		outputs[i] = &ProductionOrderOutput{
			ID:               string(order.ID),
			OrderNumber:      order.OrderNumber,
			PartID:           string(order.PartID),
			Quantity:         order.Quantity,
			Status:           string(order.Status),
			PlannedStartDate: order.Schedule.PlannedStart,
			PlannedEndDate:   order.Schedule.PlannedEnd,
			CreatedAt:        order.CreatedAt,
			UpdatedAt:        order.UpdatedAt,
		}
	}
	
	return outputs, nil
}

func (uc *ProductionUseCase) StartProduction(ctx context.Context, id string) error {
	order, err := uc.repo.FindByID(ctx, domain.ProductionOrderID(id))
	if err != nil {
		return err
	}
	
	if err := order.Start(); err != nil {
		return err
	}
	
	return uc.repo.Update(ctx, order)
}

func (uc *ProductionUseCase) CompleteProduction(ctx context.Context, id string) error {
	order, err := uc.repo.FindByID(ctx, domain.ProductionOrderID(id))
	if err != nil {
		return err
	}
	
	if err := order.Complete(); err != nil {
		return err
	}
	
	return uc.repo.Update(ctx, order)
}