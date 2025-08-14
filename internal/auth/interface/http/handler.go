package http

import (
	"database/sql"
	"encoding/json"
	"goNexttask/pkg/auth"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	db              *sql.DB
	jwtManager      *auth.JWTManager
	passwordManager *auth.PasswordManager
}

func NewAuthHandler(db *sql.DB, jwtManager *auth.JWTManager, passwordManager *auth.PasswordManager) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtManager:      jwtManager,
		passwordManager: passwordManager,
	}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/auth/register", h.Register).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST")
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.passwordManager.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	userID := uuid.New().String()
	role := req.Role
	if role == "" {
		role = "operator"
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err = h.db.Exec(query, userID, req.Email, hashedPassword, req.Name, role)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := h.jwtManager.GenerateToken(userID, req.Email, role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: UserInfo{
			ID:    userID,
			Email: req.Email,
			Name:  req.Name,
			Role:  role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var userID, hashedPassword, name, role string
	query := `SELECT id, password_hash, name, role FROM users WHERE email = $1`
	err := h.db.QueryRow(query, req.Email).Scan(&userID, &hashedPassword, &name, &role)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := h.passwordManager.VerifyPassword(hashedPassword, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtManager.GenerateToken(userID, req.Email, role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: UserInfo{
			ID:    userID,
			Email: req.Email,
			Name:  name,
			Role:  role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}