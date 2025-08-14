package domain

import (
	"context"
	"errors"
)

var (
	ErrNCProgramNotFound      = errors.New("NC program not found")
	ErrMachineNotFound        = errors.New("machine not found")
	ErrMachineNotAvailable    = errors.New("machine is not available")
	ErrIncompatibleProgram    = errors.New("program is not compatible with machine")
	ErrTransferFailed         = errors.New("program transfer failed")
)

type NCTransferService struct {
	programRepo NCProgramRepository
	machineRepo MachineRepository
}

func NewNCTransferService(programRepo NCProgramRepository, machineRepo MachineRepository) *NCTransferService {
	return &NCTransferService{
		programRepo: programRepo,
		machineRepo: machineRepo,
	}
}

func (s *NCTransferService) TransferProgram(ctx context.Context, programID NCProgramID, machineID MachineID) error {
	program, err := s.programRepo.FindByID(ctx, programID)
	if err != nil {
		return ErrNCProgramNotFound
	}
	
	machine, err := s.machineRepo.FindByID(ctx, machineID)
	if err != nil {
		return ErrMachineNotFound
	}
	
	if !machine.IsAvailable() {
		return ErrMachineNotAvailable
	}
	
	if !program.IsCompatibleWith(machine.Type) {
		return ErrIncompatibleProgram
	}
	
	// TODO: 実際のNC機器への転送ロジックを実装
	// ここではシミュレーションとして成功を返す
	
	machine.StartJob(string(programID))
	if err := s.machineRepo.Update(ctx, machine); err != nil {
		return err
	}
	
	return nil
}

func (s *NCTransferService) SelectOptimalProgram(ctx context.Context, partID string, machineType string) (*NCProgram, error) {
	// TODO: 部品IDと機械タイプから最適なNCプログラムを選定するロジック
	programs, err := s.programRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	
	for _, program := range programs {
		if program.IsCompatibleWith(machineType) {
			return program, nil
		}
	}
	
	return nil, ErrNCProgramNotFound
}