package seed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// ExtendedSeedData は製造業ナレッジフレームワークベースの拡張seedデータを投入
func ExtendedSeedData(db *sql.DB) error {
	// 1. 自動車部品製造パターン
	if err := seedAutomotivePattern(db); err != nil {
		return fmt.Errorf("failed to seed automotive pattern: %w", err)
	}
	
	// 2. 半導体製造装置パターン
	if err := seedSemiconductorPattern(db); err != nil {
		return fmt.Errorf("failed to seed semiconductor pattern: %w", err)
	}
	
	// 3. 医療機器製造パターン
	if err := seedMedicalDevicePattern(db); err != nil {
		return fmt.Errorf("failed to seed medical device pattern: %w", err)
	}
	
	// 4. 航空宇宙部品パターン
	if err := seedAerospacePattern(db); err != nil {
		return fmt.Errorf("failed to seed aerospace pattern: %w", err)
	}
	
	// 5. ロボット協働製造パターン
	if err := seedCollaborativeRobotPattern(db); err != nil {
		return fmt.Errorf("failed to seed collaborative robot pattern: %w", err)
	}
	
	return nil
}

// ====================================
// 1. 自動車部品製造パターン（制御理論統合）
// ====================================
func seedAutomotivePattern(db *sql.DB) error {
	plans := []ProductionPlan{
		// トランスミッションギア（高精度歯車）
		{
			ID:                 "PP-AUTO-2024-001",
			OrderID:            "ORD-TOYOTA-001",
			Material:           "SCM440H_歯車用鋼_浸炭焼入れ",
			Quantity:           500,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 1, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 5, 18, 0, 0, 0, time.UTC),
		},
		// エンジンクランクシャフト
		{
			ID:                 "PP-AUTO-2024-002",
			OrderID:            "ORD-HONDA-001",
			Material:           "S48C_調質材_高周波焼入れ",
			Quantity:           200,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 6, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 10, 18, 0, 0, 0, time.UTC),
		},
		// ブレーキディスクローター
		{
			ID:                 "PP-AUTO-2024-003",
			OrderID:            "ORD-NISSAN-001",
			Material:           "FC250_鋳鉄_熱処理済",
			Quantity:           1000,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 2, 20, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 28, 18, 0, 0, 0, time.UTC),
		},
		// ターボチャージャーインペラー
		{
			ID:                 "PP-AUTO-2024-004",
			OrderID:            "ORD-MAZDA-001",
			Material:           "Inconel713C_耐熱超合金",
			Quantity:           50,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 1, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 8, 18, 0, 0, 0, time.UTC),
		},
		// EV用モーターシャフト
		{
			ID:                 "PP-AUTO-2024-005",
			OrderID:            "ORD-TESLA-001",
			Material:           "SUS420J2_マルテンサイト系ステンレス",
			Quantity:           300,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 10, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 15, 18, 0, 0, 0, time.UTC),
		},
	}
	
	// NCプログラム（適応制御対応）
	ncPrograms := []NCProgram{
		{
			ID:        "NC-AUTO-GEAR-001",
			PartID:    "PART-GEAR-TMX",
			MachineID: "GLEASON-P300",
			Version:   "v5.2.0-adaptive",
			Data: `%
O5001 (TRANSMISSION GEAR - ADAPTIVE CONTROL)
(MATERIAL: SCM440H)
(PROCESS: HOBBING + SHAVING + GRINDING)
(CONTROL MODE: ADAPTIVE WITH THERMAL COMPENSATION)

G21 G90 G94 G54
M126 (HIGH PRESSURE COOLANT ON)

(ADAPTIVE FEEDRATE CONTROL)
#500=1 (ADAPTIVE MODE ON)
#501=0.85 (TARGET LOAD RATIO)
#502=100 (INITIAL FEEDRATE)

(HOBBING PROCESS WITH LOAD MONITORING)
T0101 (HOB CUTTER DIN 5480)
S800 M03
G00 X100.0 Y0 Z50.0

WHILE [#500 EQ 1] DO1
  #503=[#3002] (READ SPINDLE LOAD)
  IF [#503 GT #501] THEN #502=#502*0.95
  IF [#503 LT #501*0.8] THEN #502=#502*1.05
  G01 Z-30.0 F[#502]
  G00 Z50.0
  G91 X2.5
  G90
END1

(QUALITY PREDICTION MODEL)
#600=PREDICT_QUALITY(#503, #502, #5021)
IF [#600 LT 0.95] THEN M00 (QUALITY ALERT)

(THERMAL COMPENSATION)
#700=#5701 (MACHINE TEMPERATURE)
#701=-0.001*[#700-20] (COMPENSATION VALUE)
G10 L2 P1 Z[#5023+#701]

M30
%`,
		},
	}
	
	// 検査結果（Cpk値とPID制御フィードバック）
	inspections := []Inspection{
		{
			ID:         "INSP-AUTO-2024-001",
			LotNumber:  "LOT-GEAR-20240301-01",
			MachineID:  "CMM-ZEISS-CONTURA",
			OperatorID: "OP-AUTO-001",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"gear_pitch": {"nominal": 2.5, "actual": 2.498, "tolerance": 0.01, "unit": "mm"},
				"tooth_profile": {"deviation": 0.003, "tolerance": 0.01, "unit": "mm"},
				"helix_angle": {"nominal": 30.0, "actual": 29.98, "tolerance": 0.1, "unit": "deg"},
				"surface_hardness": {"actual": 58, "min": 56, "max": 60, "unit": "HRC"},
				"cpk": 1.72,
				"control_feedback": {
					"type": "PID",
					"error": 0.002,
					"kp": 0.5,
					"ki": 0.1,
					"kd": 0.05,
					"output": -0.001
				}
			}`),
			InspectionDate: time.Date(2024, 3, 2, 14, 30, 0, 0, time.UTC),
		},
	}
	
	// データ投入
	if err := insertProductionPlans(db, plans); err != nil {
		return err
	}
	if err := insertNCPrograms(db, ncPrograms); err != nil {
		return err
	}
	if err := insertInspections(db, inspections); err != nil {
		return err
	}
	
	return nil
}

// ====================================
// 2. 半導体製造装置パターン（超精密制御）
// ====================================
func seedSemiconductorPattern(db *sql.DB) error {
	plans := []ProductionPlan{
		// ウェハステージ部品
		{
			ID:                 "PP-SEMI-2024-001",
			OrderID:            "ORD-ASML-001",
			Material:           "SiC_セラミックス_超精密研磨",
			Quantity:           5,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 20, 17, 0, 0, 0, time.UTC),
		},
		// EUVマスクホルダー
		{
			ID:                 "PP-SEMI-2024-002",
			OrderID:            "ORD-TOKYO-ELECTRON-001",
			Material:           "Zerodur_超低膨張ガラスセラミック",
			Quantity:           2,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 21, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 4, 10, 17, 0, 0, 0, time.UTC),
		},
		// プラズマチャンバー部品
		{
			ID:                 "PP-SEMI-2024-003",
			OrderID:            "ORD-LAM-RESEARCH-001",
			Material:           "Y2O3コーティング_アルミナ",
			Quantity:           10,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 5, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 15, 17, 0, 0, 0, time.UTC),
		},
		// 真空チャックプレート
		{
			ID:                 "PP-SEMI-2024-004",
			OrderID:            "ORD-APPLIED-MATERIALS-001",
			Material:           "AlN_窒化アルミセラミック",
			Quantity:           8,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 2, 15, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 28, 17, 0, 0, 0, time.UTC),
		},
		// イオン注入装置ビームライン部品
		{
			ID:                 "PP-SEMI-2024-005",
			OrderID:            "ORD-AXCELIS-001",
			Material:           "グラファイト_高純度99.999%",
			Quantity:           15,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 25, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 4, 5, 17, 0, 0, 0, time.UTC),
		},
	}
	
	if err := insertProductionPlans(db, plans); err != nil {
		return err
	}
	
	return nil
}

// ====================================
// 3. 医療機器製造パターン（FDA/ISO13485準拠）
// ====================================
func seedMedicalDevicePattern(db *sql.DB) error {
	plans := []ProductionPlan{
		// 人工股関節ステム
		{
			ID:                 "PP-MED-2024-001",
			OrderID:            "ORD-ZIMMER-BIOMET-001",
			Material:           "Ti-6Al-4V_ELI_ASTM_F136",
			Quantity:           20,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 30, 17, 0, 0, 0, time.UTC),
		},
		// 脊椎インプラントロッド
		{
			ID:                 "PP-MED-2024-002",
			OrderID:            "ORD-MEDTRONIC-001",
			Material:           "CoCrMo_ASTM_F1537",
			Quantity:           100,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 4, 1, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 4, 15, 17, 0, 0, 0, time.UTC),
		},
		// 歯科インプラント
		{
			ID:                 "PP-MED-2024-003",
			OrderID:            "ORD-STRAUMANN-001",
			Material:           "Grade4_純チタン_SLActive処理",
			Quantity:           500,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 2, 1, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 20, 17, 0, 0, 0, time.UTC),
		},
		// 血管ステント
		{
			ID:                 "PP-MED-2024-004",
			OrderID:            "ORD-ABBOTT-001",
			Material:           "L605_コバルトクロム合金",
			Quantity:           200,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 10, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 25, 17, 0, 0, 0, time.UTC),
		},
		// 手術器具（腹腔鏡用）
		{
			ID:                 "PP-MED-2024-005",
			OrderID:            "ORD-JOHNSON-001",
			Material:           "316LVM_医療グレードステンレス",
			Quantity:           50,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 20, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 30, 17, 0, 0, 0, time.UTC),
		},
	}
	
	// 医療機器用NCプログラム（FDA準拠トレーサビリティ）
	ncPrograms := []NCProgram{
		{
			ID:        "NC-MED-HIP-001",
			PartID:    "PART-HIP-STEM",
			MachineID: "DMG-LASERTEC-65",
			Version:   "v3.0.0-FDA",
			Data: `%
O7001 (HIP IMPLANT STEM - FDA 21 CFR PART 820)
(MATERIAL: Ti-6Al-4V ELI GRADE 23)
(PROCESS: 5-AXIS MILLING + LASER TEXTURING)
(VALIDATION: IQ/OQ/PQ COMPLETED)

(TRACEABILITY HEADER)
(LOT: #LOT_NUMBER)
(OPERATOR: #OPERATOR_ID)
(DHF: DESIGN HISTORY FILE REF#12345)
(DMR: DEVICE MASTER RECORD REF#67890)

G21 G90 G94 G54.1 P1
M08 (BIOCOMPATIBLE COOLANT)

(CRITICAL DIMENSION CONTROL)
#100=12.000 (NOMINAL STEM DIAMETER)
#101=0.010 (TOLERANCE)
#102=0 (COMPENSATION)

(SURFACE TEXTURE FOR OSSEOINTEGRATION)
M160 (LASER POWER ON)
G65 P9100 A45.0 B0.050 (TEXTURE ANGLE AND DEPTH)

(IN-PROCESS VALIDATION)
G65 P9832 X[#100] Y0 Z0 T[#101]
IF [#5001 GT #101] GOTO 9999 (REJECT)

(PASSIVATION PREPARATION)
G04 P1000 (DWELL FOR SURFACE STABILIZATION)

N9999 M00 (QC HOLD POINT)
M30
%`,
		},
	}
	
	if err := insertProductionPlans(db, plans); err != nil {
		return err
	}
	if err := insertNCPrograms(db, ncPrograms); err != nil {
		return err
	}
	
	return nil
}

// ====================================
// 4. 航空宇宙部品パターン（AS9100/NADCAP準拠）
// ====================================
func seedAerospacePattern(db *sql.DB) error {
	plans := []ProductionPlan{
		// ジェットエンジンタービンブレード
		{
			ID:                 "PP-AERO-2024-001",
			OrderID:            "ORD-ROLLS-ROYCE-001",
			Material:           "CMSX-4_単結晶超合金",
			Quantity:           60,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 1, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 4, 30, 18, 0, 0, 0, time.UTC),
		},
		// 航空機主翼リブ
		{
			ID:                 "PP-AERO-2024-002",
			OrderID:            "ORD-BOEING-001",
			Material:           "Al7075-T7351_航空機構造用",
			Quantity:           20,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 5, 1, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 5, 30, 18, 0, 0, 0, time.UTC),
		},
		// ロケットエンジン燃焼室
		{
			ID:                 "PP-AERO-2024-003",
			OrderID:            "ORD-SPACEX-001",
			Material:           "GRCop-84_銅合金_3Dプリント",
			Quantity:           3,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 15, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 4, 15, 18, 0, 0, 0, time.UTC),
		},
		// 衛星用太陽電池パネル構造
		{
			ID:                 "PP-AERO-2024-004",
			OrderID:            "ORD-AIRBUS-SPACE-001",
			Material:           "CFRP_M55J_高弾性炭素繊維",
			Quantity:           8,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 2, 1, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 28, 18, 0, 0, 0, time.UTC),
		},
		// ヘリコプターメインローターハブ
		{
			ID:                 "PP-AERO-2024-005",
			OrderID:            "ORD-SIKORSKY-001",
			Material:           "Ti-10V-2Fe-3Al_チタン合金",
			Quantity:           5,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 4, 1, 7, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 4, 30, 18, 0, 0, 0, time.UTC),
		},
	}
	
	if err := insertProductionPlans(db, plans); err != nil {
		return err
	}
	
	return nil
}

// ====================================
// 5. ロボット協働製造パターン（ティーチング不要）
// ====================================
func seedCollaborativeRobotPattern(db *sql.DB) error {
	plans := []ProductionPlan{
		// 協働ロボットアーム関節
		{
			ID:                 "PP-ROBOT-2024-001",
			OrderID:            "ORD-UNIVERSAL-ROBOTS-001",
			Material:           "ADC12_精密鋳造_アルマイト処理",
			Quantity:           100,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 10, 17, 0, 0, 0, time.UTC),
		},
		// 力覚センサーハウジング
		{
			ID:                 "PP-ROBOT-2024-002",
			OrderID:            "ORD-FANUC-001",
			Material:           "A5052_アルミ合金_硬質アルマイト",
			Quantity:           200,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 11, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 20, 17, 0, 0, 0, time.UTC),
		},
		// ビジョンシステムマウント
		{
			ID:                 "PP-ROBOT-2024-003",
			OrderID:            "ORD-KUKA-001",
			Material:           "PA12_ナイロン12_SLS3Dプリント",
			Quantity:           50,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 2, 20, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 25, 17, 0, 0, 0, time.UTC),
		},
		// 安全機能付きグリッパー
		{
			ID:                 "PP-ROBOT-2024-004",
			OrderID:            "ORD-ABB-001",
			Material:           "POM_ポリアセタール_精密成形",
			Quantity:           150,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 3, 5, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 12, 17, 0, 0, 0, time.UTC),
		},
		// AIビジョン統合コントローラー筐体
		{
			ID:                 "PP-ROBOT-2024-005",
			OrderID:            "ORD-YASKAWA-001",
			Material:           "SECC_電気亜鉛メッキ鋼板",
			Quantity:           80,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 3, 25, 17, 0, 0, 0, time.UTC),
		},
	}
	
	// ティーチング不要ロボット用NCプログラム
	ncPrograms := []NCProgram{
		{
			ID:        "NC-ROBOT-ADAPTIVE-001",
			PartID:    "PART-ROBOT-JOINT",
			MachineID: "ROBOT-UR10e",
			Version:   "v2.0.0-teachless",
			Data: `%
O9001 (COLLABORATIVE ROBOT - TEACHLESS OPERATION)
(MODE: VISION-BASED AUTONOMOUS)
(SAFETY: ISO/TS 15066 COMPLIANT)

(INITIALIZE AI VISION SYSTEM)
M200 (ENABLE 3D VISION)
M201 (LOAD OBJECT RECOGNITION MODEL)

(ADAPTIVE GRIPPER CONTROL)
#800=VISION_DETECT_OBJECT()
#801=CALCULATE_GRASP_POINTS(#800)
#802=PREDICT_FORCE(#800)

(SAFETY PARAMETERS)
#900=150 (MAX FORCE IN NEWTONS)
#901=250 (MAX SPEED MM/S)
#902=1.5 (SAFETY DISTANCE METERS)

(COLLABORATIVE ZONE CHECK)
IF [HUMAN_DETECTED() EQ 1] THEN
  #901=80 (REDUCE SPEED)
  #900=50 (REDUCE FORCE)
  M203 (ENABLE SOFT SKIN)
END IF

(AUTONOMOUS PICK AND PLACE)
G00 X[#801[0]] Y[#801[1]] Z[#801[2]+100]
G01 Z[#801[2]] F[#901]
M204 (ADAPTIVE GRIP F[#802])
G00 Z[#801[2]+100]

(QUALITY LEARNING)
#950=MEASURE_SUCCESS()
M205 (UPDATE_AI_MODEL #950)

(CONTINUOUS IMPROVEMENT)
IF [#950 LT 0.95] THEN
  M206 (REQUEST_HUMAN_FEEDBACK)
END IF

M30
%`,
		},
	}
	
	// 検査結果（協働ロボット安全性検証を含む）
	inspections := []Inspection{
		{
			ID:         "INSP-ROBOT-2024-001",
			LotNumber:  "LOT-ROBOT-20240305-01",
			MachineID:  "VISION-KEYENCE-CV",
			OperatorID: "ROBOT-AI-001",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"dimensional_accuracy": {"actual": 0.02, "tolerance": 0.05, "unit": "mm"},
				"assembly_time": {"actual": 45, "target": 60, "unit": "seconds"},
				"force_applied": {"max": 120, "limit": 150, "unit": "N"},
				"safety_compliance": {
					"iso_ts_15066": "passed",
					"collision_detection": "functional",
					"speed_monitoring": "active"
				},
				"ai_performance": {
					"object_recognition_accuracy": 0.98,
					"grasp_success_rate": 0.96,
					"learning_cycles": 1250
				},
				"collaborative_metrics": {
					"human_robot_distance_min": 1.8,
					"speed_reduction_triggered": 3,
					"safety_stops": 0
				},
				"cpk": 1.89
			}`),
			InspectionDate: time.Date(2024, 3, 6, 10, 30, 0, 0, time.UTC),
		},
	}
	
	if err := insertProductionPlans(db, plans); err != nil {
		return err
	}
	if err := insertNCPrograms(db, ncPrograms); err != nil {
		return err
	}
	if err := insertInspections(db, inspections); err != nil {
		return err
	}
	
	return nil
}

// ====================================
// 共通ヘルパー関数
// ====================================
func insertProductionPlans(db *sql.DB, plans []ProductionPlan) error {
	query := `INSERT INTO production_plans 
		(id, order_id, material, quantity, status, scheduled_start_date, scheduled_end_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING`
	
	for _, plan := range plans {
		if _, err := db.Exec(query,
			plan.ID, plan.OrderID, plan.Material, plan.Quantity,
			plan.Status, plan.ScheduledStartDate, plan.ScheduledEndDate); err != nil {
			return fmt.Errorf("failed to insert production plan %s: %w", plan.ID, err)
		}
	}
	return nil
}

func insertNCPrograms(db *sql.DB, programs []NCProgram) error {
	query := `INSERT INTO nc_programs 
		(id, part_id, machine_id, version, data) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING`
	
	for _, program := range programs {
		if _, err := db.Exec(query,
			program.ID, program.PartID, program.MachineID, 
			program.Version, program.Data); err != nil {
			return fmt.Errorf("failed to insert NC program %s: %w", program.ID, err)
		}
	}
	return nil
}

func insertInspections(db *sql.DB, inspections []Inspection) error {
	query := `INSERT INTO inspections 
		(id, lot_number, machine_id, operator_id, result, measured_values, inspection_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING`
	
	for _, inspection := range inspections {
		if _, err := db.Exec(query,
			inspection.ID, inspection.LotNumber, inspection.MachineID, 
			inspection.OperatorID, inspection.Result, inspection.MeasuredValues,
			inspection.InspectionDate); err != nil {
			return fmt.Errorf("failed to insert inspection %s: %w", inspection.ID, err)
		}
	}
	return nil
}