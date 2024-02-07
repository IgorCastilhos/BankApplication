package token

import "time"

// Maker é uma interface para gerenciar tokens
type Maker interface {
	// CreateToken cria um novo token para um usuário específico com uma duração
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken verifica se o token é válido ou não
	VerifyToken(token string) (*Payload, error)
}
