package domain

import (
	"time"
)

type InspectionID string

type InspectionStatus string

const (
	InspectionStatusPending   InspectionStatus = "pending"
	InspectionStatusCompleted InspectionStatus = "completed"
	InspectionStatusFailed    InspectionStatus = "failed"
)

type InspectionResult string

const (
	ResultPass InspectionResult = "pass"
	ResultFail InspectionResult = "fail"
)

type Inspection struct {
	ID                InspectionID
	ProductionOrderID string
	LotNumber         string
	InspectorID       string
	Results           []MeasurementResult
	Status            InspectionStatus
	FinalResult       InspectionResult
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MeasurementResult struct {
	ParameterName string
	MeasuredValue float64
	TargetValue   float64
	Tolerance     float64
	Unit          string
	Pass          bool
}

type Measurement struct {
	Dimensions   map[string]float64
	InstrumentID string
	MeasuredAt   time.Time
}

func NewInspection(productionOrderID, lotNumber, inspectorID string) *Inspection {
	now := time.Now()
	return &Inspection{
		ID:                InspectionID("insp-" + time.Now().Format("20060102150405")),
		ProductionOrderID: productionOrderID,
		LotNumber:         lotNumber,
		InspectorID:       inspectorID,
		Status:            InspectionStatusPending,
		Results:           []MeasurementResult{},
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (i *Inspection) AddMeasurement(result MeasurementResult) {
	i.Results = append(i.Results, result)
	i.UpdatedAt = time.Now()
}

func (i *Inspection) Complete() {
	allPass := true
	for _, result := range i.Results {
		if !result.Pass {
			allPass = false
			break
		}
	}
	
	if allPass {
		i.FinalResult = ResultPass
	} else {
		i.FinalResult = ResultFail
	}
	
	i.Status = InspectionStatusCompleted
	i.UpdatedAt = time.Now()
}

func (i *Inspection) IsPassed() bool {
	return i.FinalResult == ResultPass
}