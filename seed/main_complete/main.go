package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"goNexttask/seed"
)

func main() {
	// コマンドラインフラグ
	var (
		reset   = flag.Bool("reset", false, "Drop all tables before setup")
		basic   = flag.Bool("basic", false, "Insert basic seed data only")
		extend  = flag.Bool("extend", true, "Insert extended seed data")
		verbose = flag.Bool("v", false, "Verbose output")
	)
	flag.Parse()

	// データベース接続情報
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "gonexttask")

	// PostgreSQL接続文字列
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// データベース接続（リトライ付き）
	var db *sql.DB
	var err error
	
	log.Println("========================================")
	log.Println("  Complete Database Setup Tool")
	log.Println("========================================")
	log.Printf("Database: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
	log.Printf("Options: reset=%v, basic=%v, extend=%v\n", *reset, *basic, *extend)
	log.Println("----------------------------------------")
	
	log.Println("Connecting to database...")
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		if *verbose {
			log.Printf("Connection attempt %d failed: %v", i+1, err)
		}
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		log.Fatalf("❌ Failed to connect to database after 30 attempts: %v", err)
	}
	defer db.Close()
	
	log.Println("✅ Connected to database successfully")
	log.Println("----------------------------------------")

	// 完全セットアップ実行
	startTime := time.Now()
	
	if *reset {
		log.Println("⚠️  WARNING: This will DROP all existing tables!")
		log.Println("Press Ctrl+C within 3 seconds to cancel...")
		time.Sleep(3 * time.Second)
	}
	
	// Complete Setup実行
	if err := seed.CompleteSetup(db); err != nil {
		log.Fatalf("❌ Failed to complete setup: %v", err)
	}
	
	duration := time.Since(startTime)
	log.Println("----------------------------------------")
	
	// 統計情報の表示
	if *verbose {
		showDetailedStatistics(db)
	} else {
		showStatistics(db)
	}
	
	log.Println("========================================")
	log.Printf("✅ Complete setup finished in %v\n", duration)
	log.Println("========================================")
}

func showStatistics(db *sql.DB) {
	fmt.Println("\n📊 Database Statistics")
	fmt.Println("----------------------")
	
	tables := []struct {
		name  string
		emoji string
	}{
		{"users", "👤"},
		{"production_plans", "📋"},
		{"production_orders", "📦"},
		{"nc_programs", "🔧"},
		{"machines", "🏭"},
		{"inspections", "🔍"},
		{"lot_inventory", "📦"},
	}
	
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table.name)
		if err := db.QueryRow(query).Scan(&count); err != nil {
			fmt.Printf("%s %-20s: ERROR\n", table.emoji, table.name)
		} else {
			fmt.Printf("%s %-20s: %d records\n", table.emoji, table.name, count)
		}
	}
}

func showDetailedStatistics(db *sql.DB) {
	showStatistics(db)
	
	fmt.Println("\n🏭 Industry Breakdown")
	fmt.Println("---------------------")
	
	industries := []struct {
		name   string
		prefix string
		emoji  string
	}{
		{"Automotive", "PP-AUTO-%", "🚗"},
		{"Semiconductor", "PP-SEMI-%", "🔬"},
		{"Medical Device", "PP-MED-%", "🏥"},
		{"Aerospace", "PP-AERO-%", "✈️"},
		{"Robotics", "PP-ROBOT-%", "🤖"},
	}
	
	for _, industry := range industries {
		var count int
		query := "SELECT COUNT(*) FROM production_plans WHERE id LIKE $1"
		if err := db.QueryRow(query, industry.prefix).Scan(&count); err != nil {
			fmt.Printf("%s %-20s: ERROR\n", industry.emoji, industry.name)
		} else {
			fmt.Printf("%s %-20s: %d plans\n", industry.emoji, industry.name, count)
		}
	}
	
	// 品質統計
	fmt.Println("\n📈 Quality Metrics")
	fmt.Println("-----------------")
	
	var totalInsp, passedInsp, failedInsp int
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN result = 'pass' THEN 1 END) as passed,
			COUNT(CASE WHEN result = 'fail' THEN 1 END) as failed
		FROM inspections
		WHERE result IS NOT NULL
	`
	
	if err := db.QueryRow(query).Scan(&totalInsp, &passedInsp, &failedInsp); err == nil && totalInsp > 0 {
		passRate := float64(passedInsp) / float64(totalInsp) * 100
		fmt.Printf("Total Inspections: %d\n", totalInsp)
		fmt.Printf("Pass Rate: %.1f%% (%d passed, %d failed)\n", passRate, passedInsp, failedInsp)
		
		// Cpk分析（JSONBから抽出）
		var avgCpk float64
		cpkQuery := `
			SELECT AVG((measured_values->>'cpk')::float)
			FROM inspections
			WHERE measured_values->>'cpk' IS NOT NULL
		`
		if err := db.QueryRow(cpkQuery).Scan(&avgCpk); err == nil {
			fmt.Printf("Average Cpk: %.2f\n", avgCpk)
		}
	}
	
	// 在庫フロー
	fmt.Println("\n📦 Inventory Flow")
	fmt.Println("-----------------")
	
	var totalIn, totalOut int
	invQuery := `
		SELECT 
			SUM(CASE WHEN in_out = 'in' THEN quantity ELSE 0 END) as total_in,
			SUM(CASE WHEN in_out = 'out' THEN quantity ELSE 0 END) as total_out
		FROM lot_inventory
	`
	
	if err := db.QueryRow(invQuery).Scan(&totalIn, &totalOut); err == nil {
		fmt.Printf("Total Inbound: %d units\n", totalIn)
		fmt.Printf("Total Outbound: %d units\n", totalOut)
		fmt.Printf("Net Inventory: %d units\n", totalIn-totalOut)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}