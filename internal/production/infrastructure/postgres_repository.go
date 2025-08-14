package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"goNexttask/internal/production/domain"
)

type PostgresProductionOrderRepository struct {
	db *sql.DB
}

func NewPostgresProductionOrderRepository(db *sql.DB) *PostgresProductionOrderRepository {
	return &PostgresProductionOrderRepository{
		db: db,
	}
}

func (r *PostgresProductionOrderRepository) Save(ctx context.Context, order *domain.ProductionOrder) error {
	machinesJSON, err := json.Marshal(order.Schedule.AssignedMachines)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO production_orders (
			id, order_number, part_id, quantity, status,
			planned_start_date, planned_end_date, assigned_machines,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = r.db.ExecContext(ctx, query,
		order.ID,
		order.OrderNumber,
		order.PartID,
		order.Quantity,
		order.Status,
		order.Schedule.PlannedStart,
		order.Schedule.PlannedEnd,
		string(machinesJSON),
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func (r *PostgresProductionOrderRepository) FindByID(ctx context.Context, id domain.ProductionOrderID) (*domain.ProductionOrder, error) {
	query := `
		SELECT id, order_number, part_id, quantity, status,
			   planned_start_date, planned_end_date, assigned_machines,
			   created_at, updated_at
		FROM production_orders
		WHERE id = $1
	`

	var order domain.ProductionOrder
	var machinesJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.PartID,
		&order.Quantity,
		&order.Status,
		&order.Schedule.PlannedStart,
		&order.Schedule.PlannedEnd,
		&machinesJSON,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrProductionOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(machinesJSON), &order.Schedule.AssignedMachines); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresProductionOrderRepository) FindAll(ctx context.Context) ([]*domain.ProductionOrder, error) {
	query := `
		SELECT id, order_number, part_id, quantity, status,
			   planned_start_date, planned_end_date, assigned_machines,
			   created_at, updated_at
		FROM production_orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.ProductionOrder

	for rows.Next() {
		var order domain.ProductionOrder
		var machinesJSON string

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.PartID,
			&order.Quantity,
			&order.Status,
			&order.Schedule.PlannedStart,
			&order.Schedule.PlannedEnd,
			&machinesJSON,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(machinesJSON), &order.Schedule.AssignedMachines); err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *PostgresProductionOrderRepository) Update(ctx context.Context, order *domain.ProductionOrder) error {
	machinesJSON, err := json.Marshal(order.Schedule.AssignedMachines)
	if err != nil {
		return err
	}

	query := `
		UPDATE production_orders
		SET order_number = $2, part_id = $3, quantity = $4, status = $5,
			planned_start_date = $6, planned_end_date = $7, assigned_machines = $8,
			updated_at = $9
		WHERE id = $1
	`

	_, err = r.db.ExecContext(ctx, query,
		order.ID,
		order.OrderNumber,
		order.PartID,
		order.Quantity,
		order.Status,
		order.Schedule.PlannedStart,
		order.Schedule.PlannedEnd,
		string(machinesJSON),
		order.UpdatedAt,
	)

	return err
}

func (r *PostgresProductionOrderRepository) Delete(ctx context.Context, id domain.ProductionOrderID) error {
	query := `DELETE FROM production_orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}