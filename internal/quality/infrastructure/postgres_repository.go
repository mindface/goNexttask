package infrastructure

import (
	"context"
	"database/sql"
	"goNexttask/internal/quality/domain"
)

type PostgresInspectionRepository struct {
	db *sql.DB
}

func NewPostgresInspectionRepository(db *sql.DB) *PostgresInspectionRepository {
	return &PostgresInspectionRepository{
		db: db,
	}
}

func (r *PostgresInspectionRepository) Save(ctx context.Context, inspection *domain.Inspection) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert inspection
	inspectionQuery := `
		INSERT INTO inspections (
			id, production_order_id, lot_number, inspector_id,
			status, final_result, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = tx.ExecContext(ctx, inspectionQuery,
		inspection.ID,
		inspection.ProductionOrderID,
		inspection.LotNumber,
		inspection.InspectorID,
		inspection.Status,
		inspection.FinalResult,
		inspection.CreatedAt,
		inspection.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert measurement results
	measurementQuery := `
		INSERT INTO measurement_results (
			inspection_id, parameter_name, measured_value,
			target_value, tolerance, unit, pass
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, result := range inspection.Results {
		_, err = tx.ExecContext(ctx, measurementQuery,
			inspection.ID,
			result.ParameterName,
			result.MeasuredValue,
			result.TargetValue,
			result.Tolerance,
			result.Unit,
			result.Pass,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresInspectionRepository) FindByID(ctx context.Context, id domain.InspectionID) (*domain.Inspection, error) {
	inspection := &domain.Inspection{}

	// Get inspection
	inspectionQuery := `
		SELECT id, production_order_id, lot_number, inspector_id,
			   status, final_result, created_at, updated_at
		FROM inspections
		WHERE id = $1
	`

	var finalResult sql.NullString
	err := r.db.QueryRowContext(ctx, inspectionQuery, id).Scan(
		&inspection.ID,
		&inspection.ProductionOrderID,
		&inspection.LotNumber,
		&inspection.InspectorID,
		&inspection.Status,
		&finalResult,
		&inspection.CreatedAt,
		&inspection.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInspectionNotFound
	}
	if err != nil {
		return nil, err
	}

	if finalResult.Valid {
		inspection.FinalResult = domain.InspectionResult(finalResult.String)
	}

	// Get measurement results
	measurementQuery := `
		SELECT parameter_name, measured_value, target_value,
			   tolerance, unit, pass
		FROM measurement_results
		WHERE inspection_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, measurementQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result domain.MeasurementResult
		err := rows.Scan(
			&result.ParameterName,
			&result.MeasuredValue,
			&result.TargetValue,
			&result.Tolerance,
			&result.Unit,
			&result.Pass,
		)
		if err != nil {
			return nil, err
		}
		inspection.Results = append(inspection.Results, result)
	}

	return inspection, nil
}

func (r *PostgresInspectionRepository) FindByLotNumber(ctx context.Context, lotNumber string) ([]*domain.Inspection, error) {
	query := `
		SELECT id, production_order_id, lot_number, inspector_id,
			   status, final_result, created_at, updated_at
		FROM inspections
		WHERE lot_number = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, lotNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []*domain.Inspection

	for rows.Next() {
		inspection := &domain.Inspection{}
		var finalResult sql.NullString

		err := rows.Scan(
			&inspection.ID,
			&inspection.ProductionOrderID,
			&inspection.LotNumber,
			&inspection.InspectorID,
			&inspection.Status,
			&finalResult,
			&inspection.CreatedAt,
			&inspection.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if finalResult.Valid {
			inspection.FinalResult = domain.InspectionResult(finalResult.String)
		}

		// Get measurement results for each inspection
		measurementQuery := `
			SELECT parameter_name, measured_value, target_value,
				   tolerance, unit, pass
			FROM measurement_results
			WHERE inspection_id = $1
			ORDER BY id
		`

		measurementRows, err := r.db.QueryContext(ctx, measurementQuery, inspection.ID)
		if err != nil {
			return nil, err
		}

		for measurementRows.Next() {
			var result domain.MeasurementResult
			err := measurementRows.Scan(
				&result.ParameterName,
				&result.MeasuredValue,
				&result.TargetValue,
				&result.Tolerance,
				&result.Unit,
				&result.Pass,
			)
			if err != nil {
				measurementRows.Close()
				return nil, err
			}
			inspection.Results = append(inspection.Results, result)
		}
		measurementRows.Close()

		inspections = append(inspections, inspection)
	}

	return inspections, nil
}

func (r *PostgresInspectionRepository) FindByProductionOrderID(ctx context.Context, orderID string) ([]*domain.Inspection, error) {
	query := `
		SELECT id, production_order_id, lot_number, inspector_id,
			   status, final_result, created_at, updated_at
		FROM inspections
		WHERE production_order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []*domain.Inspection

	for rows.Next() {
		inspection := &domain.Inspection{}
		var finalResult sql.NullString

		err := rows.Scan(
			&inspection.ID,
			&inspection.ProductionOrderID,
			&inspection.LotNumber,
			&inspection.InspectorID,
			&inspection.Status,
			&finalResult,
			&inspection.CreatedAt,
			&inspection.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if finalResult.Valid {
			inspection.FinalResult = domain.InspectionResult(finalResult.String)
		}

		inspections = append(inspections, inspection)
	}

	return inspections, nil
}

func (r *PostgresInspectionRepository) Update(ctx context.Context, inspection *domain.Inspection) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update inspection
	inspectionQuery := `
		UPDATE inspections
		SET production_order_id = $2, lot_number = $3, inspector_id = $4,
			status = $5, final_result = $6, updated_at = $7
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, inspectionQuery,
		inspection.ID,
		inspection.ProductionOrderID,
		inspection.LotNumber,
		inspection.InspectorID,
		inspection.Status,
		inspection.FinalResult,
		inspection.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Delete existing measurement results
	deleteQuery := `DELETE FROM measurement_results WHERE inspection_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, inspection.ID)
	if err != nil {
		return err
	}

	// Insert new measurement results
	measurementQuery := `
		INSERT INTO measurement_results (
			inspection_id, parameter_name, measured_value,
			target_value, tolerance, unit, pass
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, result := range inspection.Results {
		_, err = tx.ExecContext(ctx, measurementQuery,
			inspection.ID,
			result.ParameterName,
			result.MeasuredValue,
			result.TargetValue,
			result.Tolerance,
			result.Unit,
			result.Pass,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}