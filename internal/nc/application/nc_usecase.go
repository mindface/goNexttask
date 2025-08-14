package application

import (
	"context"
	"goNexttask/internal/nc/domain"
)

type RegisterNCProgramInput struct {
	Name                 string
	Version              string
	Content              []byte
	MachineCompatibility []string
	CreatedBy            string
}

type NCProgramOutput struct {
	ID                   string
	Name                 string
	Version              string
	FileHash             string
	MachineCompatibility []string
	CreatedBy            string
	CreatedAt            string
}

type DeployProgramInput struct {
	ProgramID string
	MachineID string
}

type MachineStatusOutput struct {
	ID           string
	Name         string
	IP           string
	Type         string
	RunningState string
	CurrentJobID string
	LastHeartbeat string
}

type NCUseCase struct {
	programRepo     domain.NCProgramRepository
	machineRepo     domain.MachineRepository
	transferService *domain.NCTransferService
}

func NewNCUseCase(programRepo domain.NCProgramRepository, machineRepo domain.MachineRepository) *NCUseCase {
	return &NCUseCase{
		programRepo:     programRepo,
		machineRepo:     machineRepo,
		transferService: domain.NewNCTransferService(programRepo, machineRepo),
	}
}

func (uc *NCUseCase) RegisterNCProgram(ctx context.Context, input RegisterNCProgramInput) (*NCProgramOutput, error) {
	program := domain.NewNCProgram(
		input.Name,
		input.Version,
		input.Content,
		input.MachineCompatibility,
		input.CreatedBy,
	)
	
	if err := uc.programRepo.Save(ctx, program); err != nil {
		return nil, err
	}
	
	return &NCProgramOutput{
		ID:                   string(program.ID),
		Name:                 program.Name,
		Version:              program.Version,
		FileHash:             program.FileHash,
		MachineCompatibility: program.MachineCompatibility,
		CreatedBy:            program.CreatedBy,
		CreatedAt:            program.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (uc *NCUseCase) DeployProgram(ctx context.Context, input DeployProgramInput) error {
	return uc.transferService.TransferProgram(
		ctx,
		domain.NCProgramID(input.ProgramID),
		domain.MachineID(input.MachineID),
	)
}

func (uc *NCUseCase) GetMachineStatus(ctx context.Context, machineID string) (*MachineStatusOutput, error) {
	machine, err := uc.machineRepo.FindByID(ctx, domain.MachineID(machineID))
	if err != nil {
		return nil, err
	}
	
	return &MachineStatusOutput{
		ID:            string(machine.ID),
		Name:          machine.Name,
		IP:            machine.IP,
		Type:          machine.Type,
		RunningState:  string(machine.Status.RunningState),
		CurrentJobID:  machine.Status.CurrentJobID,
		LastHeartbeat: machine.Status.LastHeartbeat.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (uc *NCUseCase) GetAllPrograms(ctx context.Context) ([]*NCProgramOutput, error) {
	programs, err := uc.programRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	
	outputs := make([]*NCProgramOutput, len(programs))
	for i, program := range programs {
		outputs[i] = &NCProgramOutput{
			ID:                   string(program.ID),
			Name:                 program.Name,
			Version:              program.Version,
			FileHash:             program.FileHash,
			MachineCompatibility: program.MachineCompatibility,
			CreatedBy:            program.CreatedBy,
			CreatedAt:            program.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	
	return outputs, nil
}

func (uc *NCUseCase) UpdateMachineStatus(ctx context.Context, machineID string, status domain.MachineStatus) error {
	machine, err := uc.machineRepo.FindByID(ctx, domain.MachineID(machineID))
	if err != nil {
		return err
	}
	
	machine.UpdateStatus(status)
	return uc.machineRepo.Update(ctx, machine)
}