package token

import "time"

// Maker é uma interface para gerenciar tokens
type Maker interface {
	// CreateToken cria um novo token para um usuário específico com uma duração
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)
	// VerifyToken verifica se o token é válido ou não
	VerifyToken(token string) (*Payload, error)
}
