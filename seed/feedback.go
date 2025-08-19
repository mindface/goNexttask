package seed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"
)

// QualityFeedbackController 品質フィードバック制御システム
type QualityFeedbackController struct {
	db              *sql.DB
	targetCpk       float64
	pidGains        PIDGains
	controlHistory  []ControlRecord
}

// PIDGains PID制御ゲイン
type PIDGains struct {
	Kp float64 // 比例ゲイン
	Ki float64 // 積分ゲイン
	Kd float64 // 微分ゲイン
}

// ControlRecord 制御履歴
type ControlRecord struct {
	Timestamp time.Time
	Error     float64
	Output    float64
	Action    string
}

// NewQualityFeedbackController 品質フィードバック制御の初期化
func NewQualityFeedbackController(db *sql.DB) *QualityFeedbackController {
	return &QualityFeedbackController{
		db:        db,
		targetCpk: 1.67, // 目標Cpk値
		pidGains: PIDGains{
			Kp: 0.5,
			Ki: 0.1,
			Kd: 0.05,
		},
		controlHistory: make([]ControlRecord, 0),
	}
}

// ProcessInspectionFeedback 検査結果からのフィードバック処理
func (qfc *QualityFeedbackController) ProcessInspectionFeedback(inspectionID string) error {
	// 検査結果の取得
	inspection, err := qfc.getInspection(inspectionID)
	if err != nil {
		return fmt.Errorf("failed to get inspection: %w", err)
	}

	// 測定値の解析
	var measuredData map[string]interface{}
	if err := json.Unmarshal(inspection.MeasuredValues, &measuredData); err != nil {
		return fmt.Errorf("failed to parse measured values: %w", err)
	}

	// Cpk値の取得
	cpk, ok := measuredData["cpk"].(float64)
	if !ok {
		cpk = qfc.calculateCpk(measuredData)
	}

	// 制御誤差の計算
	error := qfc.targetCpk - cpk
	
	// PID制御による補正値計算
	correction := qfc.calculatePIDCorrection(error)
	
	// 補正アクションの決定と実行
	if math.Abs(error) > 0.1 {
		action := qfc.determineCorrectiveAction(inspection, error, correction)
		if err := qfc.executeCorrectiveAction(action); err != nil {
			return fmt.Errorf("failed to execute corrective action: %w", err)
		}
		
		// 制御履歴の記録
		qfc.controlHistory = append(qfc.controlHistory, ControlRecord{
			Timestamp: time.Now(),
			Error:     error,
			Output:    correction,
			Action:    action.Type,
		})
	}

	return nil
}

// CorrectiveAction 補正アクション
type CorrectiveAction struct {
	Type       string
	Parameters map[string]float64
	NCProgramID string
	Description string
}

// calculateCpk Cpk値の計算
func (qfc *QualityFeedbackController) calculateCpk(measuredData map[string]interface{}) float64 {
	// 簡易的なCpk計算（実際はより複雑な統計処理が必要）
	// ここでは仮の値を返す
	return 1.5
}

// calculatePIDCorrection PID制御による補正値計算
func (qfc *QualityFeedbackController) calculatePIDCorrection(error float64) float64 {
	// 簡易PID実装
	proportional := qfc.pidGains.Kp * error
	
	// 積分項（過去のエラーの累積）
	integral := 0.0
	for _, record := range qfc.controlHistory {
		integral += record.Error
	}
	integral *= qfc.pidGains.Ki
	
	// 微分項（エラーの変化率）
	derivative := 0.0
	if len(qfc.controlHistory) > 0 {
		lastError := qfc.controlHistory[len(qfc.controlHistory)-1].Error
		derivative = qfc.pidGains.Kd * (error - lastError)
	}
	
	return proportional + integral + derivative
}

// determineCorrectiveAction 補正アクションの決定
func (qfc *QualityFeedbackController) determineCorrectiveAction(
	inspection *Inspection, 
	error float64, 
	correction float64,
) CorrectiveAction {
	action := CorrectiveAction{
		Parameters: make(map[string]float64),
	}
	
	// エラーの大きさに応じたアクション決定
	if math.Abs(error) > 0.5 {
		// 大きなエラー：NCプログラムの大幅修正
		action.Type = "major_nc_adjustment"
		action.Parameters["offset_adjustment"] = correction * 0.01 // mm単位
		action.Parameters["speed_adjustment"] = correction * -5.0  // %
		action.Description = "Major NC program adjustment due to large quality deviation"
	} else if math.Abs(error) > 0.2 {
		// 中程度のエラー：工具オフセット調整
		action.Type = "tool_offset_adjustment"
		action.Parameters["wear_offset"] = correction * 0.005 // mm単位
		action.Description = "Tool wear offset adjustment"
	} else {
		// 小さなエラー：切削条件の微調整
		action.Type = "cutting_condition_tuning"
		action.Parameters["feed_rate_adjustment"] = correction * 2.0 // %
		action.Description = "Fine tuning of cutting conditions"
	}
	
	// 失敗モードに応じた特別処理
	var measuredData map[string]interface{}
	json.Unmarshal(inspection.MeasuredValues, &measuredData)
	
	if failureMode, ok := measuredData["failure_mode"].(string); ok {
		switch failureMode {
		case "diameter_oversize":
			action.Parameters["diameter_compensation"] = -0.01
		case "surface_roughness_excess":
			action.Parameters["finish_pass_adjustment"] = 0.05
		case "concentricity_error":
			action.Parameters["alignment_correction"] = 0.005
		}
	}
	
	return action
}

// executeCorrectiveAction 補正アクションの実行
func (qfc *QualityFeedbackController) executeCorrectiveAction(action CorrectiveAction) error {
	log.Printf("Executing corrective action: %s", action.Type)
	
	switch action.Type {
	case "major_nc_adjustment":
		return qfc.updateNCProgram(action)
	case "tool_offset_adjustment":
		return qfc.updateToolOffset(action)
	case "cutting_condition_tuning":
		return qfc.updateCuttingConditions(action)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// updateNCProgram NCプログラムの更新
func (qfc *QualityFeedbackController) updateNCProgram(action CorrectiveAction) error {
	// NCプログラムの自動修正ロジック
	query := `
		UPDATE nc_programs 
		SET data = data || $1,
		    version = version || '.auto',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	
	adjustmentCode := fmt.Sprintf("\n(AUTO COMPENSATION: Offset=%.4f Speed=%.1f%%)",
		action.Parameters["offset_adjustment"],
		action.Parameters["speed_adjustment"])
	
	_, err := qfc.db.Exec(query, adjustmentCode, action.NCProgramID)
	return err
}

// updateToolOffset 工具オフセットの更新
func (qfc *QualityFeedbackController) updateToolOffset(action CorrectiveAction) error {
	// 工具オフセットテーブルへの更新（実際のシステムではマシンコントローラーへの通信）
	log.Printf("Updating tool wear offset: %.4fmm", action.Parameters["wear_offset"])
	
	// ここでは仮想的な更新を記録
	query := `
		INSERT INTO quality_adjustments (type, parameters, executed_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
	`
	
	parametersJSON, _ := json.Marshal(action.Parameters)
	_, err := qfc.db.Exec(query, action.Type, parametersJSON)
	
	return err
}

// updateCuttingConditions 切削条件の更新
func (qfc *QualityFeedbackController) updateCuttingConditions(action CorrectiveAction) error {
	log.Printf("Updating cutting conditions: Feed rate adjustment=%.1f%%", 
		action.Parameters["feed_rate_adjustment"])
	
	// 実際のシステムではCNCコントローラーへのパラメータ送信
	// ここではログ記録のみ
	return nil
}

// getInspection 検査データの取得
func (qfc *QualityFeedbackController) getInspection(inspectionID string) (*Inspection, error) {
	query := `
		SELECT id, lot_number, machine_id, operator_id, result, measured_values, inspection_date
		FROM inspections
		WHERE id = $1
	`
	
	var inspection Inspection
	err := qfc.db.QueryRow(query, inspectionID).Scan(
		&inspection.ID,
		&inspection.LotNumber,
		&inspection.MachineID,
		&inspection.OperatorID,
		&inspection.Result,
		&inspection.MeasuredValues,
		&inspection.InspectionDate,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &inspection, nil
}

// AdaptiveManufacturingSystem 適応製造システム
type AdaptiveManufacturingSystem struct {
	db                   *sql.DB
	feedbackController   *QualityFeedbackController
	learningRate         float64
	performanceHistory   []PerformanceMetric
}

// PerformanceMetric パフォーマンス指標
type PerformanceMetric struct {
	Timestamp    time.Time
	OEE          float64 // Overall Equipment Effectiveness
	QualityRate  float64
	Availability float64
	Performance  float64
}

// NewAdaptiveManufacturingSystem 適応製造システムの初期化
func NewAdaptiveManufacturingSystem(db *sql.DB) *AdaptiveManufacturingSystem {
	return &AdaptiveManufacturingSystem{
		db:                 db,
		feedbackController: NewQualityFeedbackController(db),
		learningRate:       0.01,
		performanceHistory: make([]PerformanceMetric, 0),
	}
}

// OptimizeProductionPlan 生産計画の最適化
func (ams *AdaptiveManufacturingSystem) OptimizeProductionPlan(planID string) error {
	// 過去の実績データから学習
	historicalData, err := ams.getHistoricalPerformance(planID)
	if err != nil {
		return fmt.Errorf("failed to get historical data: %w", err)
	}
	
	// モデル予測制御(MPC)による最適化
	optimalSchedule := ams.calculateOptimalSchedule(historicalData)
	
	// 生産計画の更新
	if err := ams.updateProductionPlan(planID, optimalSchedule); err != nil {
		return fmt.Errorf("failed to update production plan: %w", err)
	}
	
	return nil
}

// getHistoricalPerformance 過去のパフォーマンスデータ取得
func (ams *AdaptiveManufacturingSystem) getHistoricalPerformance(planID string) ([]PerformanceMetric, error) {
	// 実装は省略（実際はデータベースから取得）
	return ams.performanceHistory, nil
}

// calculateOptimalSchedule 最適スケジュールの計算
func (ams *AdaptiveManufacturingSystem) calculateOptimalSchedule(
	historicalData []PerformanceMetric,
) map[string]interface{} {
	// 簡易的な最適化ロジック
	schedule := make(map[string]interface{})
	
	// OEEの平均値計算
	avgOEE := 0.0
	for _, metric := range historicalData {
		avgOEE += metric.OEE
	}
	if len(historicalData) > 0 {
		avgOEE /= float64(len(historicalData))
	}
	
	// 目標OEE（85%）との差分に基づく調整
	targetOEE := 0.85
	adjustment := (targetOEE - avgOEE) * ams.learningRate
	
	schedule["cycle_time_adjustment"] = adjustment * -100 // %
	schedule["buffer_time_adjustment"] = adjustment * 50   // minutes
	schedule["batch_size_optimization"] = math.Max(10, 100*(1+adjustment))
	
	return schedule
}

// updateProductionPlan 生産計画の更新
func (ams *AdaptiveManufacturingSystem) updateProductionPlan(
	planID string,
	optimalSchedule map[string]interface{},
) error {
	// スケジュール調整の適用
	adjustmentJSON, _ := json.Marshal(optimalSchedule)
	
	query := `
		UPDATE production_plans
		SET optimization_data = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	
	_, err := ams.db.Exec(query, adjustmentJSON, planID)
	return err
}