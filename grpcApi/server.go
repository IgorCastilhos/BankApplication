package grpcApi

import (
	"fmt"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/pb"
	"github.com/IgorCastilhos/BankApplication/token"
	"github.com/IgorCastilhos/BankApplication/utils"
)

// Server define as requisições gRPC para o serviço bancário
type Server struct {
	pb.UnimplementedBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer cria um novo servidor gRPC
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar um token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
