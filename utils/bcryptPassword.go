package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword retorna o hash bcrypt da senha
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("falhou ao hashear a senha: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword verifica se a senha está correta ou não
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
