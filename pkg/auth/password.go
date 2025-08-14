package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordManager struct {
	cost int
}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		cost: bcrypt.DefaultCost,
	}
}

func (m *PasswordManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), m.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (m *PasswordManager) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}