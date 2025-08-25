package http

import (
	"encoding/json"
	"goNexttask/internal/production/application"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ProductionHandler struct {
	useCase *application.ProductionUseCase
}

func NewProductionHandler(useCase *application.ProductionUseCase) *ProductionHandler {
	return &ProductionHandler{
		useCase: useCase,
	}
}

func (h *ProductionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/production/orders", h.CreateOrder).Methods("POST")
	router.HandleFunc("/production/orders", h.GetAllOrders).Methods("GET")
	router.HandleFunc("/production/orders/{id}", h.GetOrder).Methods("GET")
	router.HandleFunc("/production/orders/{id}/start", h.StartProduction).Methods("POST")
	router.HandleFunc("/production/orders/{id}/complete", h.CompleteProduction).Methods("POST")
}

type CreateOrderRequest struct {
	OrderNumber      string    `json:"orderNumber"`
	PartID           string    `json:"partId"`
	Quantity         int       `json:"quantity"`
	PlannedStartDate time.Time `json:"plannedStartDate"`
	PlannedEndDate   time.Time `json:"plannedEndDate"`
	MachineIDs       []string  `json:"machineIds"`
}

type OrderResponse struct {
	ID               string    `json:"id"`
	OrderNumber      string    `json:"orderNumber"`
	PartID           string    `json:"partId"`
	Quantity         int       `json:"quantity"`
	Status           string    `json:"status"`
	PlannedStartDate time.Time `json:"plannedStartDate"`
	PlannedEndDate   time.Time `json:"plannedEndDate"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (h *ProductionHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := application.CreateProductionOrderInput{
		OrderNumber:      req.OrderNumber,
		PartID:           req.PartID,
		Quantity:         req.Quantity,
		PlannedStartDate: req.PlannedStartDate,
		PlannedEndDate:   req.PlannedEndDate,
		MachineIDs:       req.MachineIDs,
	}

	output, err := h.useCase.CreateProductionOrder(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := OrderResponse{
		ID:               output.ID,
		OrderNumber:      output.OrderNumber,
		PartID:           output.PartID,
		Quantity:         output.Quantity,
		Status:           output.Status,
		PlannedStartDate: output.PlannedStartDate,
		PlannedEndDate:   output.PlannedEndDate,
		CreatedAt:        output.CreatedAt,
		UpdatedAt:        output.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductionHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	output, err := h.useCase.GetProductionOrder(r.Context(), id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	response := OrderResponse{
		ID:               output.ID,
		OrderNumber:      output.OrderNumber,
		PartID:           output.PartID,
		Quantity:         output.Quantity,
		Status:           output.Status,
		PlannedStartDate: output.PlannedStartDate,
		PlannedEndDate:   output.PlannedEndDate,
		CreatedAt:        output.CreatedAt,
		UpdatedAt:        output.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductionHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	outputs, err := h.useCase.GetAllProductionOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]OrderResponse, len(outputs))
	for i, output := range outputs {
		responses[i] = OrderResponse{
			ID:               output.ID,
			OrderNumber:      output.OrderNumber,
			PartID:           output.PartID,
			Quantity:         output.Quantity,
			Status:           output.Status,
			PlannedStartDate: output.PlannedStartDate,
			PlannedEndDate:   output.PlannedEndDate,
			CreatedAt:        output.CreatedAt,
			UpdatedAt:        output.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *ProductionHandler) StartProduction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.useCase.StartProduction(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Production started"})
}

func (h *ProductionHandler) CompleteProduction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.useCase.CompleteProduction(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Production completed"})
}