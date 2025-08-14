package application

import (
	"context"
	"goNexttask/internal/quality/domain"
	"time"
)

type CreateInspectionInput struct {
	ProductionOrderID string
	LotNumber         string
	InspectorID       string
	Measurements      []MeasurementInput
}

type MeasurementInput struct {
	ParameterName string
	MeasuredValue float64
	TargetValue   float64
	Tolerance     float64
	Unit          string
}

type InspectionOutput struct {
	ID                string
	ProductionOrderID string
	LotNumber         string
	InspectorID       string
	Status            string
	FinalResult       string
	Measurements      []MeasurementOutput
	CreatedAt         time.Time
}

type MeasurementOutput struct {
	ParameterName string
	MeasuredValue float64
	TargetValue   float64
	Tolerance     float64
	Unit          string
	Pass          bool
}

type TraceabilityOutput struct {
	LotNumber   string
	Inspections []InspectionOutput
	PassRate    float64
}

type QualityUseCase struct {
	repo                  domain.InspectionRepository
	defectAnalysisService *domain.DefectAnalysisService
}

func NewQualityUseCase(repo domain.InspectionRepository) *QualityUseCase {
	return &QualityUseCase{
		repo:                  repo,
		defectAnalysisService: domain.NewDefectAnalysisService(repo),
	}
}

func (uc *QualityUseCase) CreateInspection(ctx context.Context, input CreateInspectionInput) (*InspectionOutput, error) {
	inspection := domain.NewInspection(
		input.ProductionOrderID,
		input.LotNumber,
		input.InspectorID,
	)
	
	for _, m := range input.Measurements {
		deviation := m.MeasuredValue - m.TargetValue
		pass := deviation >= -m.Tolerance && deviation <= m.Tolerance
		
		result := domain.MeasurementResult{
			ParameterName: m.ParameterName,
			MeasuredValue: m.MeasuredValue,
			TargetValue:   m.TargetValue,
			Tolerance:     m.Tolerance,
			Unit:          m.Unit,
			Pass:          pass,
		}
		
		inspection.AddMeasurement(result)
	}
	
	inspection.Complete()
	
	if err := uc.repo.Save(ctx, inspection); err != nil {
		return nil, err
	}
	
	return convertToInspectionOutput(inspection), nil
}

func (uc *QualityUseCase) GetInspection(ctx context.Context, id string) (*InspectionOutput, error) {
	inspection, err := uc.repo.FindByID(ctx, domain.InspectionID(id))
	if err != nil {
		return nil, err
	}
	
	return convertToInspectionOutput(inspection), nil
}

func (uc *QualityUseCase) GetTraceability(ctx context.Context, lotNumber string) (*TraceabilityOutput, error) {
	inspections, err := uc.repo.FindByLotNumber(ctx, lotNumber)
	if err != nil {
		return nil, err
	}
	
	analysis, err := uc.defectAnalysisService.AnalyzeDefects(ctx, lotNumber)
	if err != nil {
		return nil, err
	}
	
	outputs := make([]InspectionOutput, len(inspections))
	for i, inspection := range inspections {
		outputs[i] = *convertToInspectionOutput(inspection)
	}
	
	return &TraceabilityOutput{
		LotNumber:   lotNumber,
		Inspections: outputs,
		PassRate:    analysis.PassRate,
	}, nil
}

func (uc *QualityUseCase) AnalyzeDefects(ctx context.Context, lotNumber string) (*domain.DefectAnalysis, error) {
	return uc.defectAnalysisService.AnalyzeDefects(ctx, lotNumber)
}

func convertToInspectionOutput(inspection *domain.Inspection) *InspectionOutput {
	measurements := make([]MeasurementOutput, len(inspection.Results))
	for i, result := range inspection.Results {
		measurements[i] = MeasurementOutput{
			ParameterName: result.ParameterName,
			MeasuredValue: result.MeasuredValue,
			TargetValue:   result.TargetValue,
			Tolerance:     result.Tolerance,
			Unit:          result.Unit,
			Pass:          result.Pass,
		}
	}
	
	return &InspectionOutput{
		ID:                string(inspection.ID),
		ProductionOrderID: inspection.ProductionOrderID,
		LotNumber:         inspection.LotNumber,
		InspectorID:       inspection.InspectorID,
		Status:            string(inspection.Status),
		FinalResult:       string(inspection.FinalResult),
		Measurements:      measurements,
		CreatedAt:         inspection.CreatedAt,
	}
}