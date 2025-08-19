package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authHttp "goNexttask/internal/auth/interface/http"
	ncApp "goNexttask/internal/nc/application"
	// ncDomain "goNexttask/internal/nc/domain"
	ncHttp "goNexttask/internal/nc/interface/http"
	ncInfra "goNexttask/internal/nc/infrastructure"
	prodApp "goNexttask/internal/production/application"
	prodDomain "goNexttask/internal/production/domain"
	prodHttp "goNexttask/internal/production/interface/http"
	prodInfra "goNexttask/internal/production/infrastructure"
	qualityApp "goNexttask/internal/quality/application"
	// qualityDomain "goNexttask/internal/quality/domain"
	qualityHttp "goNexttask/internal/quality/interface/http"
	qualityInfra "goNexttask/internal/quality/infrastructure"
	"goNexttask/pkg/auth"
	"goNexttask/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "gonexttask"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database with retry
	var db *sql.DB
	var err error
	
	for i := 0; i < 30; i++ {
		db, err = database.NewConnection(dbConfig)
		if err == nil {
			break
		}
		log.Printf("Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
	}
	defer db.Close()
	
	log.Println("Database connection established")

	// Initialize auth components
	jwtManager := auth.NewJWTManager(getEnv("JWT_SECRET", "your-secret-key"), 24*time.Hour)
	passwordManager := auth.NewPasswordManager()

	// Initialize repositories
	productionRepo := prodInfra.NewPostgresProductionOrderRepository(db)
	ncProgramRepo := ncInfra.NewPostgresNCProgramRepository(db)
	machineRepo := ncInfra.NewPostgresMachineRepository(db)
	inspectionRepo := qualityInfra.NewPostgresInspectionRepository(db)

	// Initialize use cases
	productionUseCase := prodApp.NewProductionUseCase(productionRepo)
	ncUseCase := ncApp.NewNCUseCase(ncProgramRepo, machineRepo)
	qualityUseCase := qualityApp.NewQualityUseCase(inspectionRepo)

	// Initialize handlers
	authHandler := authHttp.NewAuthHandler(db, jwtManager, passwordManager)
	productionHandler := prodHttp.NewProductionHandler(productionUseCase)
	ncHandler := ncHttp.NewNCHandler(ncUseCase)
	qualityHandler := qualityHttp.NewQualityHandler(qualityUseCase)

	// Setup routes
	router := mux.NewRouter()

	// Public routes (no auth required)
	authHandler.RegisterRoutes(router)

	// Protected routes (auth required)
	protectedRouter := router.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(auth.AuthMiddleware(jwtManager))

	productionHandler.RegisterRoutes(protectedRouter)
	ncHandler.RegisterRoutes(protectedRouter)
	qualityHandler.RegisterRoutes(protectedRouter)

	// Health check endpoint with DB connection check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// DB接続確認
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(`{"status":"error","message":"Database connection failed: %v"}`, err)))
			return
		}
		
		// テーブルの存在確認
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users'").Scan(&count)
		if err != nil || count == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(`{"status":"error","message":"Users table not found: %v"}`, err)))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","database":"connected","tables":"exists"}`))
	}).Methods("GET")

	// Setup server
	srv := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Stub implementations for missing repositories
// These need to be implemented in their respective infrastructure packages

type stubProductionRepo struct{}

func (r *stubProductionRepo) Save(ctx context.Context, order *prodDomain.ProductionOrder) error {
	return nil
}
func (r *stubProductionRepo) FindByID(ctx context.Context, id prodDomain.ProductionOrderID) (*prodDomain.ProductionOrder, error) {
	return nil, fmt.Errorf("not implemented")
}
func (r *stubProductionRepo) FindAll(ctx context.Context) ([]*prodDomain.ProductionOrder, error) {
	return nil, nil
}
func (r *stubProductionRepo) Update(ctx context.Context, order *prodDomain.ProductionOrder) error {
	return nil
}
func (r *stubProductionRepo) Delete(ctx context.Context, id prodDomain.ProductionOrderID) error {
	return nil
}