package grpcApi

import (
	"context"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/pb"
	"github.com/IgorCastilhos/BankApplication/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "falhou ao fazer o hash da senha: %s", err)
	}

	// Se for válido...
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "nome de usuário já existe: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "falhou ao criar um usuário: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}
