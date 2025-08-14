package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type NCProgramID string

type NCProgram struct {
	ID                   NCProgramID
	Name                 string
	Version              string
	FileHash             string
	MachineCompatibility []string
	CreatedBy            string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewNCProgram(name, version string, content []byte, machineCompatibility []string, createdBy string) *NCProgram {
	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])
	
	now := time.Now()
	return &NCProgram{
		ID:                   NCProgramID("ncprog-" + hashStr[:8]),
		Name:                 name,
		Version:              version,
		FileHash:             hashStr,
		MachineCompatibility: machineCompatibility,
		CreatedBy:            createdBy,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

func (p *NCProgram) IsCompatibleWith(machineType string) bool {
	for _, mt := range p.MachineCompatibility {
		if mt == machineType {
			return true
		}
	}
	return false
}