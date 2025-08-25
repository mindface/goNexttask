package seed

import (
	"database/sql"
	"fmt"
	"log"
)

// CompleteSetup はデータベースの完全セットアップ（削除→作成→seed）を実行
func CompleteSetup(db *sql.DB) error {
	log.Println("Starting complete database setup...")
	
	// 1. 既存テーブルの削除
	if err := dropAllTables(db); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	
	// 2. スキーマの作成
	if err := createAllTables(db); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	
	// 3. 基本seedデータの投入
	if err := SeedData(db); err != nil {
		return fmt.Errorf("failed to seed basic data: %w", err)
	}
	
	// 4. 拡張seedデータの投入
	if err := ExtendedSeedData(db); err != nil {
		return fmt.Errorf("failed to seed extended data: %w", err)
	}
	
	log.Println("Complete database setup finished successfully!")
	return nil
}

// dropAllTables は全テーブルを削除
func dropAllTables(db *sql.DB) error {
	log.Println("Dropping all tables...")
	
	tables := []string{
		"measurement_results",  // 外部キー依存があるため先に削除
		"inspections",
		"lot_inventory",
		"production_plans",
		"production_orders",
		"nc_programs",
		"machines",
		"users",
		"schema_migrations",
		"quality_adjustments",  // feedback.goで使用される可能性
	}
	
	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to drop table %s: %v", table, err)
			// エラーがあっても続行（テーブルが存在しない可能性）
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}
	
	return nil
}

// createAllTables は全テーブルを作成
func createAllTables(db *sql.DB) error {
	log.Println("Creating all tables...")
	
	// スキーマ管理テーブル
	if err := createSchemaMigrationsTable(db); err != nil {
		return err
	}
	
	// ユーザー認証テーブル
	if err := createUsersTable(db); err != nil {
		return err
	}
	
	// 生産管理テーブル
	if err := createProductionTables(db); err != nil {
		return err
	}
	
	// NC加工連携テーブル
	if err := createNCTables(db); err != nil {
		return err
	}
	
	// 品質管理テーブル
	if err := createQualityTables(db); err != nil {
		return err
	}
	
	// 在庫管理テーブル
	if err := createInventoryTables(db); err != nil {
		return err
	}
	
	// 拡張テーブル（フィードバック用）
	if err := createExtensionTables(db); err != nil {
		return err
	}
	
	// インデックスの作成
	if err := createIndexes(db); err != nil {
		return err
	}
	
	log.Println("All tables created successfully!")
	return nil
}

func createSchemaMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}
	log.Println("Created table: schema_migrations")
	return nil
}

func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(64) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		name VARCHAR(128) NOT NULL,
		role VARCHAR(32) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Println("Created table: users")
	return nil
}

func createProductionTables(db *sql.DB) error {
	// production_orders (既存のDDD設計用)
	query1 := `
	CREATE TABLE IF NOT EXISTS production_orders (
		id VARCHAR(64) PRIMARY KEY,
		order_number VARCHAR(128) NOT NULL UNIQUE,
		part_id VARCHAR(64) NOT NULL,
		quantity INT NOT NULL,
		status VARCHAR(32) NOT NULL CHECK (status IN ('planned', 'in_progress', 'completed', 'delayed', 'cancelled')),
		planned_start_date TIMESTAMP NOT NULL,
		planned_end_date TIMESTAMP NOT NULL,
		assigned_machines TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query1); err != nil {
		return fmt.Errorf("failed to create production_orders: %w", err)
	}
	log.Println("Created table: production_orders")
	
	// production_plans (seedデータ用)
	query2 := `
	CREATE TABLE IF NOT EXISTS production_plans (
		id VARCHAR(64) PRIMARY KEY,
		order_id VARCHAR(64) NOT NULL,
		material VARCHAR(256) NOT NULL,
		quantity INT NOT NULL,
		status VARCHAR(32) NOT NULL CHECK (status IN ('planned', 'in_progress', 'completed', 'delayed')),
		scheduled_start_date TIMESTAMP NOT NULL,
		scheduled_end_date TIMESTAMP NOT NULL,
		optimization_data JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query2); err != nil {
		return fmt.Errorf("failed to create production_plans: %w", err)
	}
	log.Println("Created table: production_plans")
	return nil
}

func createNCTables(db *sql.DB) error {
	// NCプログラム管理
	query1 := `
	CREATE TABLE IF NOT EXISTS nc_programs (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(128),
		part_id VARCHAR(64),
		machine_id VARCHAR(64),
		version VARCHAR(32) NOT NULL,
		file_hash VARCHAR(256),
		machine_compatibility TEXT,
		data TEXT NOT NULL,
		created_by VARCHAR(128),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query1); err != nil {
		return fmt.Errorf("failed to create nc_programs: %w", err)
	}
	log.Println("Created table: nc_programs")
	
	// 機械管理
	query2 := `
	CREATE TABLE IF NOT EXISTS machines (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(128) NOT NULL,
		ip_address VARCHAR(45) NOT NULL,
		machine_type VARCHAR(64) NOT NULL,
		capabilities TEXT,
		running_state VARCHAR(32) NOT NULL CHECK (running_state IN ('running', 'stopped', 'error')),
		current_job_id VARCHAR(64),
		last_heartbeat TIMESTAMP,
		error_message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query2); err != nil {
		return fmt.Errorf("failed to create machines: %w", err)
	}
	log.Println("Created table: machines")
	return nil
}

func createQualityTables(db *sql.DB) error {
	// 検査結果
	query1 := `
	CREATE TABLE IF NOT EXISTS inspections (
		id VARCHAR(64) PRIMARY KEY,
		production_order_id VARCHAR(64),
		lot_number VARCHAR(64) NOT NULL,
		machine_id VARCHAR(64),
		inspector_id VARCHAR(64),
		operator_id VARCHAR(64),
		status VARCHAR(32) CHECK (status IN ('pending', 'completed', 'failed')),
		result VARCHAR(16) CHECK (result IN ('pass', 'fail')),
		final_result VARCHAR(16) CHECK (final_result IN ('pass', 'fail')),
		measured_values JSONB,
		inspection_date TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query1); err != nil {
		return fmt.Errorf("failed to create inspections: %w", err)
	}
	log.Println("Created table: inspections")
	
	// 測定結果詳細
	query2 := `
	CREATE TABLE IF NOT EXISTS measurement_results (
		id SERIAL PRIMARY KEY,
		inspection_id VARCHAR(64) NOT NULL REFERENCES inspections(id) ON DELETE CASCADE,
		parameter_name VARCHAR(128) NOT NULL,
		measured_value DECIMAL(10, 4) NOT NULL,
		target_value DECIMAL(10, 4) NOT NULL,
		tolerance DECIMAL(10, 4) NOT NULL,
		unit VARCHAR(32) NOT NULL,
		pass BOOLEAN NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query2); err != nil {
		return fmt.Errorf("failed to create measurement_results: %w", err)
	}
	log.Println("Created table: measurement_results")
	return nil
}

func createInventoryTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS lot_inventory (
		id VARCHAR(64) PRIMARY KEY,
		lot_number VARCHAR(64) NOT NULL,
		product_type VARCHAR(128) NOT NULL,
		quantity INT NOT NULL,
		in_out VARCHAR(8) NOT NULL CHECK (in_out IN ('in', 'out')),
		transaction_date TIMESTAMP NOT NULL,
		location VARCHAR(128),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create lot_inventory: %w", err)
	}
	log.Println("Created table: lot_inventory")
	return nil
}

func createExtensionTables(db *sql.DB) error {
	// 品質調整履歴（feedback.goで使用）
	query1 := `
	CREATE TABLE IF NOT EXISTS quality_adjustments (
		id SERIAL PRIMARY KEY,
		type VARCHAR(64) NOT NULL,
		parameters JSONB,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query1); err != nil {
		return fmt.Errorf("failed to create quality_adjustments: %w", err)
	}
	log.Println("Created table: quality_adjustments")
	
	// 購入注文（在庫管理用）
	query2 := `
	CREATE TABLE IF NOT EXISTS purchase_orders (
		id SERIAL PRIMARY KEY,
		product_type VARCHAR(128) NOT NULL,
		quantity INT NOT NULL,
		status VARCHAR(32) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query2); err != nil {
		return fmt.Errorf("failed to create purchase_orders: %w", err)
	}
	log.Println("Created table: purchase_orders")
	
	// 品質アラート
	query3 := `
	CREATE TABLE IF NOT EXISTS quality_alerts (
		id SERIAL PRIMARY KEY,
		inspection_id VARCHAR(64),
		alert_type VARCHAR(64) NOT NULL,
		message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := db.Exec(query3); err != nil {
		return fmt.Errorf("failed to create quality_alerts: %w", err)
	}
	log.Println("Created table: quality_alerts")
	
	return nil
}

func createIndexes(db *sql.DB) error {
	log.Println("Creating indexes...")
	
	indexes := []string{
		// production_orders
		"CREATE INDEX IF NOT EXISTS idx_production_orders_status ON production_orders(status)",
		"CREATE INDEX IF NOT EXISTS idx_production_orders_dates ON production_orders(planned_start_date, planned_end_date)",
		
		// production_plans
		"CREATE INDEX IF NOT EXISTS idx_production_plans_status ON production_plans(status)",
		"CREATE INDEX IF NOT EXISTS idx_production_plans_order ON production_plans(order_id)",
		"CREATE INDEX IF NOT EXISTS idx_production_plans_dates ON production_plans(scheduled_start_date, scheduled_end_date)",
		
		// machines
		"CREATE INDEX IF NOT EXISTS idx_machines_state ON machines(running_state)",
		
		// inspections
		"CREATE INDEX IF NOT EXISTS idx_inspections_lot ON inspections(lot_number)",
		"CREATE INDEX IF NOT EXISTS idx_inspections_order ON inspections(production_order_id)",
		"CREATE INDEX IF NOT EXISTS idx_inspections_result ON inspections(result)",
		"CREATE INDEX IF NOT EXISTS idx_inspections_date ON inspections(inspection_date)",
		
		// lot_inventory
		"CREATE INDEX IF NOT EXISTS idx_lot_inventory_lot ON lot_inventory(lot_number)",
		"CREATE INDEX IF NOT EXISTS idx_lot_inventory_product ON lot_inventory(product_type)",
		"CREATE INDEX IF NOT EXISTS idx_lot_inventory_date ON lot_inventory(transaction_date)",
		
		// users
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)",
	}
	
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
			// インデックス作成エラーは警告のみ（続行）
		}
	}
	
	log.Println("Indexes created successfully!")
	return nil
}