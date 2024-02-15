package grpcApi

import (
	"context"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/pb"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/IgorCastilhos/BankApplication/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
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

func validateCreateUserRequest(request *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidatePassword(request.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validation.ValidateFullName(request.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validation.ValidateEmail(request.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
