package domain

import (
	"context"
	"errors"
)

var (
	ErrInspectionNotFound = errors.New("inspection not found")
	ErrInvalidMeasurement = errors.New("invalid measurement")
)

type DefectAnalysisService struct {
	repo InspectionRepository
}

func NewDefectAnalysisService(repo InspectionRepository) *DefectAnalysisService {
	return &DefectAnalysisService{
		repo: repo,
	}
}

func (s *DefectAnalysisService) AnalyzeDefects(ctx context.Context, lotNumber string) (*DefectAnalysis, error) {
	inspections, err := s.repo.FindByLotNumber(ctx, lotNumber)
	if err != nil {
		return nil, err
	}
	
	if len(inspections) == 0 {
		return nil, ErrInspectionNotFound
	}
	
	analysis := &DefectAnalysis{
		LotNumber:     lotNumber,
		TotalSamples:  len(inspections),
		PassedSamples: 0,
		FailedSamples: 0,
		DefectTypes:   make(map[string]int),
	}
	
	for _, inspection := range inspections {
		if inspection.IsPassed() {
			analysis.PassedSamples++
		} else {
			analysis.FailedSamples++
			for _, result := range inspection.Results {
				if !result.Pass {
					analysis.DefectTypes[result.ParameterName]++
				}
			}
		}
	}
	
	analysis.PassRate = float64(analysis.PassedSamples) / float64(analysis.TotalSamples) * 100
	
	return analysis, nil
}

type DefectAnalysis struct {
	LotNumber     string
	TotalSamples  int
	PassedSamples int
	FailedSamples int
	PassRate      float64
	DefectTypes   map[string]int
}

func (s *DefectAnalysisService) GetTraceability(ctx context.Context, lotNumber string) (*TraceabilityInfo, error) {
	inspections, err := s.repo.FindByLotNumber(ctx, lotNumber)
	if err != nil {
		return nil, err
	}
	
	if len(inspections) == 0 {
		return nil, ErrInspectionNotFound
	}
	
	info := &TraceabilityInfo{
		LotNumber:   lotNumber,
		Inspections: inspections,
		// TODO: NCプログラムバージョン、工具情報、機械ログとの関連付け
	}
	
	return info, nil
}

type TraceabilityInfo struct {
	LotNumber        string
	Inspections      []*Inspection
	NCProgramVersion string
	MachineID        string
	ToolInfo         map[string]interface{}
}