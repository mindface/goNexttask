package domain

import "context"

type NCProgramRepository interface {
	Save(ctx context.Context, program *NCProgram) error
	FindByID(ctx context.Context, id NCProgramID) (*NCProgram, error)
	FindByNameAndVersion(ctx context.Context, name, version string) (*NCProgram, error)
	FindAll(ctx context.Context) ([]*NCProgram, error)
	Delete(ctx context.Context, id NCProgramID) error
}

type MachineRepository interface {
	Save(ctx context.Context, machine *Machine) error
	FindByID(ctx context.Context, id MachineID) (*Machine, error)
	FindAll(ctx context.Context) ([]*Machine, error)
	FindAvailable(ctx context.Context) ([]*Machine, error)
	Update(ctx context.Context, machine *Machine) error
}