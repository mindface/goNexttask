package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"goNexttask/internal/nc/domain"
	"time"
)

type PostgresNCProgramRepository struct {
	db *sql.DB
}

func NewPostgresNCProgramRepository(db *sql.DB) *PostgresNCProgramRepository {
	return &PostgresNCProgramRepository{
		db: db,
	}
}

func (r *PostgresNCProgramRepository) Save(ctx context.Context, program *domain.NCProgram) error {
	compatibilityJSON, err := json.Marshal(program.MachineCompatibility)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO nc_programs (id, name, version, file_hash, machine_compatibility, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.ExecContext(ctx, query,
		program.ID,
		program.Name,
		program.Version,
		program.FileHash,
		string(compatibilityJSON),
		program.CreatedBy,
		program.CreatedAt,
		program.UpdatedAt,
	)

	return err
}

func (r *PostgresNCProgramRepository) FindByID(ctx context.Context, id domain.NCProgramID) (*domain.NCProgram, error) {
	query := `
		SELECT id, name, version, file_hash, machine_compatibility, created_by, created_at, updated_at
		FROM nc_programs
		WHERE id = $1
	`

	var program domain.NCProgram
	var compatibilityJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&program.ID,
		&program.Name,
		&program.Version,
		&program.FileHash,
		&compatibilityJSON,
		&program.CreatedBy,
		&program.CreatedAt,
		&program.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNCProgramNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(compatibilityJSON), &program.MachineCompatibility); err != nil {
		return nil, err
	}

	return &program, nil
}

func (r *PostgresNCProgramRepository) FindByNameAndVersion(ctx context.Context, name, version string) (*domain.NCProgram, error) {
	query := `
		SELECT id, name, version, file_hash, machine_compatibility, created_by, created_at, updated_at
		FROM nc_programs
		WHERE name = $1 AND version = $2
	`

	var program domain.NCProgram
	var compatibilityJSON string

	err := r.db.QueryRowContext(ctx, query, name, version).Scan(
		&program.ID,
		&program.Name,
		&program.Version,
		&program.FileHash,
		&compatibilityJSON,
		&program.CreatedBy,
		&program.CreatedAt,
		&program.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNCProgramNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(compatibilityJSON), &program.MachineCompatibility); err != nil {
		return nil, err
	}

	return &program, nil
}

func (r *PostgresNCProgramRepository) FindAll(ctx context.Context) ([]*domain.NCProgram, error) {
	query := `
		SELECT id, name, version, file_hash, machine_compatibility, created_by, created_at, updated_at
		FROM nc_programs
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []*domain.NCProgram

	for rows.Next() {
		var program domain.NCProgram
		var compatibilityJSON string

		err := rows.Scan(
			&program.ID,
			&program.Name,
			&program.Version,
			&program.FileHash,
			&compatibilityJSON,
			&program.CreatedBy,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(compatibilityJSON), &program.MachineCompatibility); err != nil {
			return nil, err
		}

		programs = append(programs, &program)
	}

	return programs, nil
}

func (r *PostgresNCProgramRepository) Delete(ctx context.Context, id domain.NCProgramID) error {
	query := `DELETE FROM nc_programs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

type PostgresMachineRepository struct {
	db *sql.DB
}

func NewPostgresMachineRepository(db *sql.DB) *PostgresMachineRepository {
	return &PostgresMachineRepository{
		db: db,
	}
}

func (r *PostgresMachineRepository) Save(ctx context.Context, machine *domain.Machine) error {
	capabilitiesJSON, err := json.Marshal(machine.Capabilities)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO machines (
			id, name, ip_address, machine_type, capabilities,
			running_state, current_job_id, last_heartbeat, error_message,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.db.ExecContext(ctx, query,
		machine.ID,
		machine.Name,
		machine.IP,
		machine.Type,
		string(capabilitiesJSON),
		machine.Status.RunningState,
		machine.Status.CurrentJobID,
		machine.Status.LastHeartbeat,
		machine.Status.ErrorMessage,
		machine.CreatedAt,
		machine.UpdatedAt,
	)

	return err
}

func (r *PostgresMachineRepository) FindByID(ctx context.Context, id domain.MachineID) (*domain.Machine, error) {
	query := `
		SELECT id, name, ip_address, machine_type, capabilities,
			   running_state, current_job_id, last_heartbeat, error_message,
			   created_at, updated_at
		FROM machines
		WHERE id = $1
	`

	var machine domain.Machine
	var capabilitiesJSON string
	var currentJobID sql.NullString
	var errorMessage sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&machine.ID,
		&machine.Name,
		&machine.IP,
		&machine.Type,
		&capabilitiesJSON,
		&machine.Status.RunningState,
		&currentJobID,
		&machine.Status.LastHeartbeat,
		&errorMessage,
		&machine.CreatedAt,
		&machine.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrMachineNotFound
	}
	if err != nil {
		return nil, err
	}

	if currentJobID.Valid {
		machine.Status.CurrentJobID = currentJobID.String
	}
	if errorMessage.Valid {
		machine.Status.ErrorMessage = errorMessage.String
	}

	if err := json.Unmarshal([]byte(capabilitiesJSON), &machine.Capabilities); err != nil {
		return nil, err
	}

	return &machine, nil
}

func (r *PostgresMachineRepository) FindAll(ctx context.Context) ([]*domain.Machine, error) {
	query := `
		SELECT id, name, ip_address, machine_type, capabilities,
			   running_state, current_job_id, last_heartbeat, error_message,
			   created_at, updated_at
		FROM machines
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []*domain.Machine

	for rows.Next() {
		var machine domain.Machine
		var capabilitiesJSON string
		var currentJobID sql.NullString
		var errorMessage sql.NullString

		err := rows.Scan(
			&machine.ID,
			&machine.Name,
			&machine.IP,
			&machine.Type,
			&capabilitiesJSON,
			&machine.Status.RunningState,
			&currentJobID,
			&machine.Status.LastHeartbeat,
			&errorMessage,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if currentJobID.Valid {
			machine.Status.CurrentJobID = currentJobID.String
		}
		if errorMessage.Valid {
			machine.Status.ErrorMessage = errorMessage.String
		}

		if err := json.Unmarshal([]byte(capabilitiesJSON), &machine.Capabilities); err != nil {
			return nil, err
		}

		machines = append(machines, &machine)
	}

	return machines, nil
}

func (r *PostgresMachineRepository) FindAvailable(ctx context.Context) ([]*domain.Machine, error) {
	query := `
		SELECT id, name, ip_address, machine_type, capabilities,
			   running_state, current_job_id, last_heartbeat, error_message,
			   created_at, updated_at
		FROM machines
		WHERE running_state = 'stopped'
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []*domain.Machine

	for rows.Next() {
		var machine domain.Machine
		var capabilitiesJSON string
		var currentJobID sql.NullString
		var errorMessage sql.NullString

		err := rows.Scan(
			&machine.ID,
			&machine.Name,
			&machine.IP,
			&machine.Type,
			&capabilitiesJSON,
			&machine.Status.RunningState,
			&currentJobID,
			&machine.Status.LastHeartbeat,
			&errorMessage,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if currentJobID.Valid {
			machine.Status.CurrentJobID = currentJobID.String
		}
		if errorMessage.Valid {
			machine.Status.ErrorMessage = errorMessage.String
		}

		if err := json.Unmarshal([]byte(capabilitiesJSON), &machine.Capabilities); err != nil {
			return nil, err
		}

		machines = append(machines, &machine)
	}

	return machines, nil
}

func (r *PostgresMachineRepository) Update(ctx context.Context, machine *domain.Machine) error {
	capabilitiesJSON, err := json.Marshal(machine.Capabilities)
	if err != nil {
		return err
	}

	query := `
		UPDATE machines
		SET name = $2, ip_address = $3, machine_type = $4, capabilities = $5,
			running_state = $6, current_job_id = $7, last_heartbeat = $8, error_message = $9,
			updated_at = $10
		WHERE id = $1
	`

	_, err = r.db.ExecContext(ctx, query,
		machine.ID,
		machine.Name,
		machine.IP,
		machine.Type,
		string(capabilitiesJSON),
		machine.Status.RunningState,
		machine.Status.CurrentJobID,
		machine.Status.LastHeartbeat,
		machine.Status.ErrorMessage,
		time.Now(),
	)

	return err
}