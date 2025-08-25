package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"goNexttask/seed"
)

func main() {
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
	
	log.Println("Connecting to database...")
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
	}
	defer db.Close()
	
	log.Println("Connected to database successfully")

	// 拡張シードデータ投入（30パターン）
	log.Println("Inserting extended seed data (30 patterns)...")
	if err := seed.ExtendedSeedData(db); err != nil {
		log.Fatalf("Failed to seed extended data: %v", err)
	}

	// 統計情報の表示
	showStatistics(db)
	
	log.Println("Extended seed data inserted successfully!")
}

func showStatistics(db *sql.DB) {
	type TableCount struct {
		TableName string
		Count     int
	}
	
	tables := []string{"production_plans", "nc_programs", "inspections", "lot_inventory"}
	
	fmt.Println("\n=== Seed Data Statistics ===")
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		if err := db.QueryRow(query).Scan(&count); err != nil {
			log.Printf("Error counting %s: %v", table, err)
			continue
		}
		fmt.Printf("%-20s: %d records\n", table, count)
	}
	
	// 業界別の統計
	fmt.Println("\n=== Industry Pattern Statistics ===")
	industryPrefixes := map[string]string{
		"Automotive":     "PP-AUTO-%",
		"Semiconductor":  "PP-SEMI-%",
		"Medical Device": "PP-MED-%",
		"Aerospace":      "PP-AERO-%",
		"Robotics":       "PP-ROBOT-%",
	}
	
	for industry, prefix := range industryPrefixes {
		var count int
		query := "SELECT COUNT(*) FROM production_plans WHERE id LIKE $1"
		if err := db.QueryRow(query, prefix).Scan(&count); err != nil {
			log.Printf("Error counting %s: %v", industry, err)
			continue
		}
		fmt.Printf("%-20s: %d plans\n", industry, count)
	}
	
	// 品質統計（Cpk値の分布）
	fmt.Println("\n=== Quality Statistics (Cpk Distribution) ===")
	query := `
		SELECT 
			COUNT(*) as total_inspections,
			COUNT(CASE WHEN result = 'pass' THEN 1 END) as passed,
			COUNT(CASE WHEN result = 'fail' THEN 1 END) as failed
		FROM inspections
	`
	var total, passed, failed int
	if err := db.QueryRow(query).Scan(&total, &passed, &failed); err == nil {
		fmt.Printf("Total Inspections: %d\n", total)
		fmt.Printf("Passed: %d (%.1f%%)\n", passed, float64(passed)/float64(total)*100)
		fmt.Printf("Failed: %d (%.1f%%)\n", failed, float64(failed)/float64(total)*100)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}