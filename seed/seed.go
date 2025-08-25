package seed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// SeedData はデータベースに初期データを投入する
func SeedData(db *sql.DB) error {
	if err := seedProductionPlans(db); err != nil {
		return fmt.Errorf("failed to seed production_plans: %w", err)
	}
	
	if err := seedNCPrograms(db); err != nil {
		return fmt.Errorf("failed to seed nc_programs: %w", err)
	}
	
	if err := seedInspections(db); err != nil {
		return fmt.Errorf("failed to seed inspections: %w", err)
	}
	
	if err := seedLotInventory(db); err != nil {
		return fmt.Errorf("failed to seed lot_inventory: %w", err)
	}
	
	return nil
}

// ProductionPlan represents a production plan record
type ProductionPlan struct {
	ID                string
	OrderID           string
	Material          string
	Quantity          int
	Status            string
	ScheduledStartDate time.Time
	ScheduledEndDate   time.Time
}

func seedProductionPlans(db *sql.DB) error {
	plans := []ProductionPlan{
		// 自動車部品
		{
			ID:                 "PP-2024-001",
			OrderID:            "ORD-AUTO-001",
			Material:           "SCM440H_φ50x300_調質",
			Quantity:           100,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 1, 20, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "PP-2024-002",
			OrderID:            "ORD-AUTO-002",
			Material:           "S45C_φ30x200_焼入れ",
			Quantity:           200,
			Status:             "completed",
			ScheduledStartDate: time.Date(2024, 1, 18, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 1, 25, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "PP-2024-003",
			OrderID:            "ORD-AUTO-003",
			Material:           "SUS304_φ25x150_研磨",
			Quantity:           150,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 1, 22, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 1, 28, 17, 0, 0, 0, time.UTC),
		},
		// 半導体製造装置部品
		{
			ID:                 "PP-2024-004",
			OrderID:            "ORD-SEMI-001",
			Material:           "A7075-T6_200x150x50_精密加工",
			Quantity:           20,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 1, 23, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 1, 30, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "PP-2024-005",
			OrderID:            "ORD-SEMI-002",
			Material:           "SUS316L_φ100x30_鏡面仕上げ",
			Quantity:           30,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 1, 29, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 5, 17, 0, 0, 0, time.UTC),
		},
		// 医療機器部品
		{
			ID:                 "PP-2024-006",
			OrderID:            "ORD-MED-001",
			Material:           "Ti-6Al-4V_人工関節_カスタム",
			Quantity:           5,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 1, 24, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 10, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "PP-2024-007",
			OrderID:            "ORD-MED-002",
			Material:           "PEEK_インプラント_φ8x20",
			Quantity:           50,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 2, 1, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 8, 17, 0, 0, 0, time.UTC),
		},
		// 航空宇宙部品
		{
			ID:                 "PP-2024-008",
			OrderID:            "ORD-AERO-001",
			Material:           "Inconel718_タービンブレード",
			Quantity:           10,
			Status:             "delayed",
			ScheduledStartDate: time.Date(2024, 1, 20, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 15, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:                 "PP-2024-009",
			OrderID:            "ORD-AERO-002",
			Material:           "CFRP_主翼部品_1200x800x50",
			Quantity:           2,
			Status:             "planned",
			ScheduledStartDate: time.Date(2024, 2, 5, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 28, 17, 0, 0, 0, time.UTC),
		},
		// ロボット部品 
		{
			ID:                 "PP-2024-010",
			OrderID:            "ORD-ROBO-001",
			Material:           "ADC12_ロボットアーム筐体",
			Quantity:           15,
			Status:             "in_progress",
			ScheduledStartDate: time.Date(2024, 1, 25, 8, 0, 0, 0, time.UTC),
			ScheduledEndDate:   time.Date(2024, 2, 2, 17, 0, 0, 0, time.UTC),
		},
	}

	query := `INSERT INTO production_plans 
		(id, order_id, material, quantity, status, scheduled_start_date, scheduled_end_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, plan := range plans {
		if _, err := db.Exec(query,
			plan.ID, plan.OrderID, plan.Material, plan.Quantity,
			plan.Status, plan.ScheduledStartDate, plan.ScheduledEndDate); err != nil {
			return fmt.Errorf("failed to insert production plan %s: %w", plan.ID, err)
		}
	}
	
	return nil
}

// NCProgram represents an NC program record
type NCProgram struct {
	ID        string
	PartID    string
	MachineID string
	Version   string
	Data      string
}

func seedNCPrograms(db *sql.DB) error {
	programs := []NCProgram{
		{
			ID:        "NC-SHAFT-001",
			PartID:    "PART-AUTO-001",
			MachineID: "DMG-MORI-NTX2000",
			Version:   "v3.2.1",
			Data: `%
O0001 (HIGH PRECISION SHAFT)
(MATERIAL: SCM440H DIA50X300)
(TOOL LIST:)
(T01: ROUGH TURNING INSERT CNMG120408)
(T02: FINISH TURNING INSERT DNMG110404)
(T03: GROOVING TOOL 3MM)
(T04: THREADING TOOL 60DEG)

G21 G90 G94 G54
M08 (COOLANT ON)
G96 S180 M03 (CONSTANT SURFACE SPEED)

(ROUGH TURNING)
T0101
G00 X52.0 Z2.0
G71 U1.0 R0.5
G71 P100 Q200 U0.5 W0.1 F0.25
N100 G00 X48.0
G01 Z-295.0 F0.2
X50.0 Z-298.0
N200 X52.0

(FINISH TURNING WITH ADAPTIVE CONTROL)
T0202
G00 X50.5 Z2.0
G70 P100 Q200
G04 P500 (DWELL FOR THERMAL STABILIZATION)

(GROOVING OPERATIONS)
T0303
G00 X52.0 Z-50.0
G75 R0.5
G75 X46.0 Z-50.0 P1000 Q500 F0.05

(QUALITY CHECK PROBE)
T0909
G65 P9832 (AUTOMATIC MEASUREMENT MACRO)
M30
%`,
		},
		{
			ID:        "NC-SEMI-001",
			PartID:    "PART-SEMI-001",
			MachineID: "MAKINO-A51NX",
			Version:   "v2.8.3",
			Data: `%
O0002 (SEMICONDUCTOR EQUIPMENT PART)
(MATERIAL: A7075-T6 200X150X50)
(TOLERANCE: ±0.002MM)
(SURFACE FINISH: RA0.1)

G21 G90 G94 G54
M08 (MINIMUM QUANTITY LUBRICATION)

(HIGH SPEED MACHINING STRATEGY)
G05.1 Q1 (AI SERVO ON)
G43.4 (TOOL CENTER POINT CONTROL)

(ADAPTIVE FEEDRATE OPTIMIZATION)
#100=120 (INITIAL FEEDRATE)
#101=0 (TOOL WEAR COMPENSATION)
WHILE [#101 LT 0.05] DO1
  G01 X[#102] Y[#103] F[#100-#101*50]
  #101=#5001 (READ SPINDLE LOAD)
END1

(THERMAL COMPENSATION)
G10 L2 P1 X[#5021+#500] Y[#5022+#501] Z[#5023+#502]

(IN-PROCESS MEASUREMENT)
G65 P9810 A1.0 B0.002 (TOLERANCE CHECK MACRO)
M30
%`,
		},
		{
			ID:        "NC-MED-001",
			PartID:    "PART-MED-001",
			MachineID: "HERMLE-C400",
			Version:   "v4.1.0",
			Data: `%
O0003 (MEDICAL IMPLANT - HIP JOINT)
(MATERIAL: TI-6AL-4V GRADE 5)
(FDA 21 CFR PART 820 COMPLIANT)
(FULL TRACEABILITY ENABLED)

G21 G90 G94 G54.1 P1
M08 (FLOOD COOLANT - MEDICAL GRADE)

(TROCHOIDAL MILLING FOR HEAT MANAGEMENT)
G05 P10000 (LOOK AHEAD 10000 BLOCKS)
G61.1 (HIGH PRECISION MODE)

(CONSTANT CHIP LOAD STRATEGY)
#1=0.05 (TARGET CHIP THICKNESS)
#2=8000 (SPINDLE RPM)
#3=#2*#1*4 (FEEDRATE CALCULATION)

(5-AXIS SIMULTANEOUS)
G43.5 H01 (TCP WITH KINEMATICS)
G68.2 P1 Q1 R1 (TILTED WORK PLANE)

(SURFACE GENERATION WITH NURBS)
G06.2 K1 P3 (NURBS INTERPOLATION MODE)
N100 X10.5 Y15.3 Z-5.2 A30.0 B45.0 F[#3]
N110 X11.2 Y16.1 Z-5.5 A30.5 B45.2

(ULTRASONIC ASSISTED MACHINING)
M160 (ULTRASONIC ON)
G01 X20.0 Y25.0 F300
M161 (ULTRASONIC OFF)

M30
%`,
		},
		{
			ID:        "NC-AERO-001",
			PartID:    "PART-AERO-001",
			MachineID: "MAZAK-INTEGREX-I400",
			Version:   "v5.0.2",
			Data: `%
O0004 (TURBINE BLADE - INCONEL 718)
(AS9100D COMPLIANT)
(NADCAP CERTIFIED PROCESS)

G21 G90 G94 G54
M126 (HIGH PRESSURE COOLANT 70BAR)

(CERAMIC TOOL STRATEGY)
T0101 (CERAMIC INSERT)
S450 M03
G96 S60 (CONSTANT SURFACE SPEED FOR CERAMICS)

(INTELLIGENT TOOL LIFE MANAGEMENT)
IF [#3001 GT 1800] GOTO 9999 (TOOL LIFE CHECK)

(PLUNGE MILLING FOR DIFFICULT MATERIALS)
G00 X50.0 Y0 Z5.0
#10=0 (PLUNGE COUNTER)
WHILE [#10 LT 100] DO2
  G01 Z-30.0 F50
  G00 Z5.0
  G91 G01 X2.0
  G90
  #10=#10+1
END2

(CRYOGENIC COOLING ACTIVATION)
M240 (LN2 COOLING ON)
G04 P2000 (STABILIZATION)

(BLADE PROFILE WITH COMPENSATION)
G41 D01 (CUTTER COMP LEFT)
G01 X[#5041] Y[#5042] F80
G40 (CANCEL CUTTER COMP)

N9999 M01 (OPTIONAL STOP FOR TOOL CHANGE)
M30
%`,
		},
		{
			ID:        "NC-ROBO-001",
			PartID:    "PART-ROBO-001",
			MachineID: "BROTHER-SPEEDIO",
			Version:   "v3.5.0",
			Data: `%
O0005 (ROBOT ARM HOUSING)
(MATERIAL: ADC12 ALUMINUM DIE CAST)
(COLLABORATIVE ROBOT COMPONENT)

G21 G90 G94 G54
M08 (COOLANT ON)

(HIGH SPEED MACHINING)
G05.1 Q1 R5 (HPCC MODE)
S20000 M03

(TOOL PATH OPTIMIZATION BY AI)
M200 (ENABLE AI OPTIMIZATION)
#500=1 (LEARNING MODE ON)

(DYNAMIC FIXTURE COMPENSATION)
G10 L2 P1 X[#5221] Y[#5222] Z[#5223]

(POCKET MILLING WITH CHIP EVACUATION)
T0202 (6MM ENDMILL)
M06
G43 H02
G00 X10.0 Y10.0 Z5.0
G01 Z-10.0 F2000

(HELICAL INTERPOLATION)
G02 X20.0 Y10.0 I5.0 J0 Z-15.0 F1500

(AUTOMATIC DEBURRING PATH)
T0808 (DEBURRING TOOL)
M06
G65 P9850 (DEBURRING MACRO)

(COORDINATE MEASURING)
T0909
G65 P9832 X10.0 Y10.0 Z0 A0.1 (MEASURE AND COMPENSATE)

M30
%`,
		},
	}

	query := `INSERT INTO nc_programs (id, part_id, machine_id, version, data) VALUES ($1, $2, $3, $4, $5)`

	for _, program := range programs {
		if _, err := db.Exec(query,
			program.ID, program.PartID, program.MachineID, program.Version, program.Data); err != nil {
			return fmt.Errorf("failed to insert NC program %s: %w", program.ID, err)
		}
	}

	return nil
}

// Inspection represents an inspection record
type Inspection struct {
	ID             string
	LotNumber      string
	MachineID      string
	OperatorID     string
	Result         string
	MeasuredValues json.RawMessage
	InspectionDate time.Time
}

func seedInspections(db *sql.DB) error {
	inspections := []Inspection{
		{
			ID:         "INSP-2024-001",
			LotNumber:  "LOT-AUTO-20240115-01",
			MachineID:  "CMM-ZEISS-01",
			OperatorID: "OP-001",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"diameter": {"nominal": 50.0, "actual": 49.998, "tolerance": 0.01, "unit": "mm"},
				"length": {"nominal": 300.0, "actual": 299.999, "tolerance": 0.05, "unit": "mm"},
				"concentricity": {"actual": 0.003, "tolerance": 0.01, "unit": "mm"},
				"surface_roughness": {"actual": 0.4, "tolerance": 0.8, "unit": "Ra"},
				"hardness": {"actual": 58, "tolerance_min": 56, "tolerance_max": 60, "unit": "HRC"},
				"cpk": 1.67,
				"temperature_during_measurement": 20.0
			}`),
			InspectionDate: time.Date(2024, 1, 20, 14, 30, 0, 0, time.UTC),
		},
		{
			ID:         "INSP-2024-002",
			LotNumber:  "LOT-AUTO-20240118-01",
			MachineID:  "CMM-ZEISS-01",
			OperatorID: "OP-002",
			Result:     "fail",
			MeasuredValues: json.RawMessage(`{
				"diameter": {"nominal": 30.0, "actual": 30.015, "tolerance": 0.01, "unit": "mm"},
				"length": {"nominal": 200.0, "actual": 200.002, "tolerance": 0.05, "unit": "mm"},
				"concentricity": {"actual": 0.012, "tolerance": 0.01, "unit": "mm"},
				"surface_roughness": {"actual": 0.6, "tolerance": 0.8, "unit": "Ra"},
				"failure_mode": "diameter_oversize",
				"root_cause": "tool_wear",
				"corrective_action": "nc_program_offset_adjustment",
				"cpk": 0.89
			}`),
			InspectionDate: time.Date(2024, 1, 25, 10, 15, 0, 0, time.UTC),
		},
		{
			ID:         "INSP-2024-003",
			LotNumber:  "LOT-SEMI-20240123-01",
			MachineID:  "CMM-ZEISS-PRISMO",
			OperatorID: "OP-003",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"flatness": {"actual": 0.0008, "tolerance": 0.002, "unit": "mm"},
				"parallelism": {"actual": 0.0012, "tolerance": 0.003, "unit": "mm"},
				"position_tolerance": {"actual": 0.0015, "tolerance": 0.005, "unit": "mm"},
				"surface_roughness": {"actual": 0.05, "tolerance": 0.1, "unit": "Ra"},
				"cleanliness": {"particle_count": 12, "max_allowed": 50, "class": "ISO14644-1_Class5"},
				"measurement_uncertainty": 0.0002,
				"cpk": 2.01
			}`),
			InspectionDate: time.Date(2024, 1, 26, 16, 45, 0, 0, time.UTC),
		},
		{
			ID:         "INSP-2024-004",
			LotNumber:  "LOT-MED-20240124-01",
			MachineID:  "CMM-HEXAGON-MEDICAL",
			OperatorID: "OP-004",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"dimensional_accuracy": {"actual": 0.003, "tolerance": 0.01, "unit": "mm"},
				"surface_finish": {"actual": 0.08, "tolerance": 0.2, "unit": "Ra"},
				"biocompatibility": "passed",
				"sterilization_validation": "passed",
				"fatigue_test": {"cycles": 1000000, "result": "no_failure"},
				"traceability_id": "FDA-2024-0124-001",
				"cpk": 1.85
			}`),
			InspectionDate: time.Date(2024, 1, 30, 9, 20, 0, 0, time.UTC),
		},
		{
			ID:         "INSP-2024-005",
			LotNumber:  "LOT-AERO-20240120-01",
			MachineID:  "CMM-LEITZ-PMM",
			OperatorID: "OP-005",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"blade_profile": {"deviation": 0.008, "tolerance": 0.02, "unit": "mm"},
				"surface_integrity": {"cracks": "none", "residual_stress": -120, "unit": "MPa"},
				"material_certification": "AMS5663",
				"heat_treatment": "solution_treated_aged",
				"ndt_inspection": {"method": "FPI", "result": "no_defects"},
				"far_compliance": "FAR_Part_21",
				"cpk": 1.72
			}`),
			InspectionDate: time.Date(2024, 2, 10, 11, 30, 0, 0, time.UTC),
		},
		{
			ID:         "INSP-2024-006",
			LotNumber:  "LOT-ROBO-20240125-01",
			MachineID:  "CMM-MITUTOYO-01",
			OperatorID: "OP-006",
			Result:     "pass",
			MeasuredValues: json.RawMessage(`{
				"mounting_hole_position": {"actual": 0.02, "tolerance": 0.05, "unit": "mm"},
				"weight": {"actual": 1.245, "nominal": 1.250, "tolerance": 0.050, "unit": "kg"},
				"surface_quality": {"defects": "none", "coating_thickness": 25, "unit": "μm"},
				"iso_ts_15066_compliance": "passed",
				"force_limitation_test": {"max_force": 140, "limit": 150, "unit": "N"},
				"cpk": 1.45
			}`),
			InspectionDate: time.Date(2024, 1, 28, 13, 50, 0, 0, time.UTC),
		},
	}

	query := `INSERT INTO inspections 
		(id, lot_number, machine_id, operator_id, result, measured_values, inspection_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, inspection := range inspections {
		if _, err := db.Exec(query,
			inspection.ID, inspection.LotNumber, inspection.MachineID, inspection.OperatorID,
			inspection.Result, inspection.MeasuredValues, inspection.InspectionDate); err != nil {
			return fmt.Errorf("failed to insert inspection %s: %w", inspection.ID, err)
		}
	}
	
	return nil
}

// LotInventory represents a lot inventory record
type LotInventory struct {
	ID              string
	LotNumber       string
	ProductType     string
	Quantity        int
	InOut           string
	TransactionDate time.Time
}

func seedLotInventory(db *sql.DB) error {
	inventories := []LotInventory{
		// 原材料入庫
		{
			ID:              "INV-2024-001",
			LotNumber:       "MAT-SCM440H-20240110",
			ProductType:     "原材料_SCM440H_φ50",
			Quantity:        500,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 10, 8, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-002",
			LotNumber:       "MAT-S45C-20240112",
			ProductType:     "原材料_S45C_φ30",
			Quantity:        1000,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 12, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-003",
			LotNumber:       "MAT-SUS304-20240113",
			ProductType:     "原材料_SUS304_φ25",
			Quantity:        800,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 13, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-004",
			LotNumber:       "MAT-A7075-20240114",
			ProductType:     "原材料_A7075-T6_ブロック",
			Quantity:        100,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 14, 11, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-005",
			LotNumber:       "MAT-TI6AL4V-20240115",
			ProductType:     "原材料_Ti-6Al-4V_ブロック",
			Quantity:        20,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 15, 8, 30, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-006",
			LotNumber:       "MAT-INCONEL-20240116",
			ProductType:     "原材料_Inconel718_棒材",
			Quantity:        50,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 16, 9, 30, 0, 0, time.UTC),
		},
		// 半製品（工程間移動）
		{
			ID:              "INV-2024-007",
			LotNumber:       "LOT-AUTO-20240115-01",
			ProductType:     "半製品_シャフト_旋削完了",
			Quantity:        100,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 17, 15, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-008",
			LotNumber:       "LOT-AUTO-20240115-01",
			ProductType:     "半製品_シャフト_旋削完了",
			Quantity:        100,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 18, 8, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-009",
			LotNumber:       "LOT-AUTO-20240115-01",
			ProductType:     "半製品_シャフト_熱処理完了",
			Quantity:        100,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 18, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-010",
			LotNumber:       "LOT-AUTO-20240115-01",
			ProductType:     "半製品_シャフト_熱処理完了",
			Quantity:        100,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 19, 8, 0, 0, 0, time.UTC),
		},
		// 完成品出庫
		{
			ID:              "INV-2024-011",
			LotNumber:       "LOT-AUTO-20240115-01",
			ProductType:     "完成品_高精度シャフト",
			Quantity:        98,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 20, 15, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-012",
			LotNumber:       "LOT-AUTO-20240115-01",
			ProductType:     "完成品_高精度シャフト",
			Quantity:        98,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 21, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-013",
			LotNumber:       "LOT-AUTO-20240118-01",
			ProductType:     "完成品_精密シャフト",
			Quantity:        195,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 25, 16, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-014",
			LotNumber:       "LOT-AUTO-20240118-01",
			ProductType:     "完成品_精密シャフト",
			Quantity:        195,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 26, 9, 0, 0, 0, time.UTC),
		},
		// 不良品隔離
		{
			ID:              "INV-2024-015",
			LotNumber:       "LOT-AUTO-20240118-01",
			ProductType:     "不良品_径寸法超過",
			Quantity:        5,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 25, 10, 30, 0, 0, time.UTC),
		},
		// JIT納入（かんばん方式）
		{
			ID:              "INV-2024-016",
			LotNumber:       "LOT-SEMI-20240123-01",
			ProductType:     "完成品_半導体装置部品",
			Quantity:        20,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 26, 17, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-017",
			LotNumber:       "LOT-SEMI-20240123-01",
			ProductType:     "完成品_半導体装置部品",
			Quantity:        10,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 27, 8, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-018",
			LotNumber:       "LOT-SEMI-20240123-01",
			ProductType:     "完成品_半導体装置部品",
			Quantity:        10,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 28, 8, 0, 0, 0, time.UTC),
		},
		// 医療機器（厳格なロット管理）
		{
			ID:              "INV-2024-019",
			LotNumber:       "LOT-MED-20240124-01",
			ProductType:     "完成品_人工関節_滅菌済",
			Quantity:        5,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 30, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-020",
			LotNumber:       "LOT-MED-20240124-01",
			ProductType:     "完成品_人工関節_滅菌済",
			Quantity:        1,
			InOut:           "out",
			TransactionDate: time.Date(2024, 1, 31, 14, 0, 0, 0, time.UTC),
		},
		// リサイクル材料
		{
			ID:              "INV-2024-021",
			LotNumber:       "RECYCLE-AL-20240125",
			ProductType:     "リサイクル_アルミ切粉",
			Quantity:        50,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 25, 16, 0, 0, 0, time.UTC),
		},
		{
			ID:              "INV-2024-022",
			LotNumber:       "RECYCLE-SUS-20240126",
			ProductType:     "リサイクル_SUS切粉",
			Quantity:        30,
			InOut:           "in",
			TransactionDate: time.Date(2024, 1, 26, 16, 0, 0, 0, time.UTC),
		},
	}

	query := `INSERT INTO lot_inventory 
		(id, lot_number, product_type, quantity, in_out, transaction_date) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	for _, inventory := range inventories {
		if _, err := db.Exec(query,
			inventory.ID, inventory.LotNumber, inventory.ProductType,
			inventory.Quantity, inventory.InOut, inventory.TransactionDate); err != nil {
			return fmt.Errorf("failed to insert lot inventory %s: %w", inventory.ID, err)
		}
	}
	
	return nil
}