package http

import (
	"encoding/json"
	"goNexttask/internal/quality/application"
	"net/http"

	"github.com/gorilla/mux"
)

type QualityHandler struct {
	useCase *application.QualityUseCase
}

func NewQualityHandler(useCase *application.QualityUseCase) *QualityHandler {
	return &QualityHandler{
		useCase: useCase,
	}
}

func (h *QualityHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/quality/inspections", h.CreateInspection).Methods("POST")
	router.HandleFunc("/quality/inspections/{id}", h.GetInspection).Methods("GET")
	router.HandleFunc("/quality/traceability", h.GetTraceability).Methods("GET")
	router.HandleFunc("/quality/defect-analysis", h.AnalyzeDefects).Methods("GET")
}

type CreateInspectionRequest struct {
	ProductionOrderID string                  `json:"productionOrderId"`
	LotNumber         string                  `json:"lotNumber"`
	InspectorID       string                  `json:"inspectorId"`
	Measurements      []MeasurementRequest    `json:"measurements"`
}

type MeasurementRequest struct {
	ParameterName string  `json:"parameterName"`
	MeasuredValue float64 `json:"measuredValue"`
	TargetValue   float64 `json:"targetValue"`
	Tolerance     float64 `json:"tolerance"`
	Unit          string  `json:"unit"`
}

type InspectionResponse struct {
	ID                string                 `json:"id"`
	ProductionOrderID string                 `json:"productionOrderId"`
	LotNumber         string                 `json:"lotNumber"`
	InspectorID       string                 `json:"inspectorId"`
	Status            string                 `json:"status"`
	FinalResult       string                 `json:"finalResult"`
	Measurements      []MeasurementResponse  `json:"measurements"`
	CreatedAt         string                 `json:"createdAt"`
}

type MeasurementResponse struct {
	ParameterName string  `json:"parameterName"`
	MeasuredValue float64 `json:"measuredValue"`
	TargetValue   float64 `json:"targetValue"`
	Tolerance     float64 `json:"tolerance"`
	Unit          string  `json:"unit"`
	Pass          bool    `json:"pass"`
}

type TraceabilityResponse struct {
	LotNumber   string               `json:"lotNumber"`
	Inspections []InspectionResponse `json:"inspections"`
	PassRate    float64              `json:"passRate"`
}

type DefectAnalysisResponse struct {
	LotNumber     string         `json:"lotNumber"`
	TotalSamples  int            `json:"totalSamples"`
	PassedSamples int            `json:"passedSamples"`
	FailedSamples int            `json:"failedSamples"`
	PassRate      float64        `json:"passRate"`
	DefectTypes   map[string]int `json:"defectTypes"`
}

func (h *QualityHandler) CreateInspection(w http.ResponseWriter, r *http.Request) {
	var req CreateInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	measurements := make([]application.MeasurementInput, len(req.Measurements))
	for i, m := range req.Measurements {
		measurements[i] = application.MeasurementInput{
			ParameterName: m.ParameterName,
			MeasuredValue: m.MeasuredValue,
			TargetValue:   m.TargetValue,
			Tolerance:     m.Tolerance,
			Unit:          m.Unit,
		}
	}

	input := application.CreateInspectionInput{
		ProductionOrderID: req.ProductionOrderID,
		LotNumber:         req.LotNumber,
		InspectorID:       req.InspectorID,
		Measurements:      measurements,
	}

	output, err := h.useCase.CreateInspection(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	measurementResponses := make([]MeasurementResponse, len(output.Measurements))
	for i, m := range output.Measurements {
		measurementResponses[i] = MeasurementResponse{
			ParameterName: m.ParameterName,
			MeasuredValue: m.MeasuredValue,
			TargetValue:   m.TargetValue,
			Tolerance:     m.Tolerance,
			Unit:          m.Unit,
			Pass:          m.Pass,
		}
	}

	response := InspectionResponse{
		ID:                output.ID,
		ProductionOrderID: output.ProductionOrderID,
		LotNumber:         output.LotNumber,
		InspectorID:       output.InspectorID,
		Status:            output.Status,
		FinalResult:       output.FinalResult,
		Measurements:      measurementResponses,
		CreatedAt:         output.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *QualityHandler) GetInspection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	output, err := h.useCase.GetInspection(r.Context(), id)
	if err != nil {
		http.Error(w, "Inspection not found", http.StatusNotFound)
		return
	}

	measurementResponses := make([]MeasurementResponse, len(output.Measurements))
	for i, m := range output.Measurements {
		measurementResponses[i] = MeasurementResponse{
			ParameterName: m.ParameterName,
			MeasuredValue: m.MeasuredValue,
			TargetValue:   m.TargetValue,
			Tolerance:     m.Tolerance,
			Unit:          m.Unit,
			Pass:          m.Pass,
		}
	}

	response := InspectionResponse{
		ID:                output.ID,
		ProductionOrderID: output.ProductionOrderID,
		LotNumber:         output.LotNumber,
		InspectorID:       output.InspectorID,
		Status:            output.Status,
		FinalResult:       output.FinalResult,
		Measurements:      measurementResponses,
		CreatedAt:         output.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *QualityHandler) GetTraceability(w http.ResponseWriter, r *http.Request) {
	lotNumber := r.URL.Query().Get("lot")
	if lotNumber == "" {
		http.Error(w, "Lot number is required", http.StatusBadRequest)
		return
	}

	output, err := h.useCase.GetTraceability(r.Context(), lotNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inspectionResponses := make([]InspectionResponse, len(output.Inspections))
	for i, insp := range output.Inspections {
		measurementResponses := make([]MeasurementResponse, len(insp.Measurements))
		for j, m := range insp.Measurements {
			measurementResponses[j] = MeasurementResponse{
				ParameterName: m.ParameterName,
				MeasuredValue: m.MeasuredValue,
				TargetValue:   m.TargetValue,
				Tolerance:     m.Tolerance,
				Unit:          m.Unit,
				Pass:          m.Pass,
			}
		}

		inspectionResponses[i] = InspectionResponse{
			ID:                insp.ID,
			ProductionOrderID: insp.ProductionOrderID,
			LotNumber:         insp.LotNumber,
			InspectorID:       insp.InspectorID,
			Status:            insp.Status,
			FinalResult:       insp.FinalResult,
			Measurements:      measurementResponses,
			CreatedAt:         insp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	response := TraceabilityResponse{
		LotNumber:   output.LotNumber,
		Inspections: inspectionResponses,
		PassRate:    output.PassRate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *QualityHandler) AnalyzeDefects(w http.ResponseWriter, r *http.Request) {
	lotNumber := r.URL.Query().Get("lot")
	if lotNumber == "" {
		http.Error(w, "Lot number is required", http.StatusBadRequest)
		return
	}

	analysis, err := h.useCase.AnalyzeDefects(r.Context(), lotNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DefectAnalysisResponse{
		LotNumber:     analysis.LotNumber,
		TotalSamples:  analysis.TotalSamples,
		PassedSamples: analysis.PassedSamples,
		FailedSamples: analysis.FailedSamples,
		PassRate:      analysis.PassRate,
		DefectTypes:   analysis.DefectTypes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}