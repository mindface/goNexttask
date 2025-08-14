package http

import (
	"encoding/json"
	"goNexttask/internal/nc/application"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type NCHandler struct {
	useCase *application.NCUseCase
}

func NewNCHandler(useCase *application.NCUseCase) *NCHandler {
	return &NCHandler{
		useCase: useCase,
	}
}

func (h *NCHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/nc/programs", h.RegisterProgram).Methods("POST")
	router.HandleFunc("/api/v1/nc/programs", h.GetAllPrograms).Methods("GET")
	router.HandleFunc("/api/v1/nc/machines/{id}/deploy", h.DeployProgram).Methods("POST")
	router.HandleFunc("/api/v1/nc/machines/{id}/status", h.GetMachineStatus).Methods("GET")
	router.HandleFunc("/api/v1/nc/machines/{id}/status", h.UpdateMachineStatus).Methods("POST")
}

type RegisterProgramRequest struct {
	Name                 string   `json:"name"`
	Version              string   `json:"version"`
	Content              string   `json:"content"`
	MachineCompatibility []string `json:"machineCompatibility"`
	CreatedBy            string   `json:"createdBy"`
}

type ProgramResponse struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Version              string   `json:"version"`
	FileHash             string   `json:"fileHash"`
	MachineCompatibility []string `json:"machineCompatibility"`
	CreatedBy            string   `json:"createdBy"`
	CreatedAt            string   `json:"createdAt"`
}

type DeployRequest struct {
	ProgramID string `json:"programId"`
}

type MachineStatusRequest struct {
	RunningState string `json:"runningState"`
	CurrentJobID string `json:"currentJobId,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type MachineStatusResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IP            string `json:"ip"`
	Type          string `json:"type"`
	RunningState  string `json:"runningState"`
	CurrentJobID  string `json:"currentJobId"`
	LastHeartbeat string `json:"lastHeartbeat"`
}

func (h *NCHandler) RegisterProgram(w http.ResponseWriter, r *http.Request) {
	var req RegisterProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := application.RegisterNCProgramInput{
		Name:                 req.Name,
		Version:              req.Version,
		Content:              []byte(req.Content),
		MachineCompatibility: req.MachineCompatibility,
		CreatedBy:            req.CreatedBy,
	}

	output, err := h.useCase.RegisterNCProgram(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ProgramResponse{
		ID:                   output.ID,
		Name:                 output.Name,
		Version:              output.Version,
		FileHash:             output.FileHash,
		MachineCompatibility: output.MachineCompatibility,
		CreatedBy:            output.CreatedBy,
		CreatedAt:            output.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *NCHandler) GetAllPrograms(w http.ResponseWriter, r *http.Request) {
	outputs, err := h.useCase.GetAllPrograms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]ProgramResponse, len(outputs))
	for i, output := range outputs {
		responses[i] = ProgramResponse{
			ID:                   output.ID,
			Name:                 output.Name,
			Version:              output.Version,
			FileHash:             output.FileHash,
			MachineCompatibility: output.MachineCompatibility,
			CreatedBy:            output.CreatedBy,
			CreatedAt:            output.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *NCHandler) DeployProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineID := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req DeployRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := application.DeployProgramInput{
		ProgramID: req.ProgramID,
		MachineID: machineID,
	}

	if err := h.useCase.DeployProgram(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Program deployed successfully"})
}

func (h *NCHandler) GetMachineStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineID := vars["id"]

	output, err := h.useCase.GetMachineStatus(r.Context(), machineID)
	if err != nil {
		http.Error(w, "Machine not found", http.StatusNotFound)
		return
	}

	response := MachineStatusResponse{
		ID:            output.ID,
		Name:          output.Name,
		IP:            output.IP,
		Type:          output.Type,
		RunningState:  output.RunningState,
		CurrentJobID:  output.CurrentJobID,
		LastHeartbeat: output.LastHeartbeat,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *NCHandler) UpdateMachineStatus(w http.ResponseWriter, r *http.Request) {
	var req MachineStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Convert request to domain.MachineStatus and update
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Machine status updated"})
}