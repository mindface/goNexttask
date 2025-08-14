package domain

import "time"

type MachineID string

type MachineRunningState string

const (
	StateRunning MachineRunningState = "running"
	StateStopped MachineRunningState = "stopped"
	StateError   MachineRunningState = "error"
)

type Machine struct {
	ID           MachineID
	Name         string
	IP           string
	Type         string
	Capabilities []string
	Status       MachineStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MachineStatus struct {
	RunningState  MachineRunningState
	CurrentJobID  string
	LastHeartbeat time.Time
	ErrorMessage  string
}

func NewMachine(id MachineID, name, ip, machineType string, capabilities []string) *Machine {
	now := time.Now()
	return &Machine{
		ID:           id,
		Name:         name,
		IP:           ip,
		Type:         machineType,
		Capabilities: capabilities,
		Status: MachineStatus{
			RunningState:  StateStopped,
			LastHeartbeat: now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (m *Machine) UpdateStatus(status MachineStatus) {
	m.Status = status
	m.UpdatedAt = time.Now()
}

func (m *Machine) IsAvailable() bool {
	return m.Status.RunningState == StateStopped
}

func (m *Machine) StartJob(jobID string) {
	m.Status.RunningState = StateRunning
	m.Status.CurrentJobID = jobID
	m.Status.LastHeartbeat = time.Now()
	m.UpdatedAt = time.Now()
}

func (m *Machine) StopJob() {
	m.Status.RunningState = StateStopped
	m.Status.CurrentJobID = ""
	m.Status.LastHeartbeat = time.Now()
	m.UpdatedAt = time.Now()
}