package grpcApi

import (
	"context"
	"errors"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/pb"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/IgorCastilhos/BankApplication/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "usuário não encontrado")
		}
		return nil, status.Errorf(codes.Internal, "falha ao buscar usuário")
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "senha incorreta")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username, user.Role, server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "falha ao criar token de acesso")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "falha ao criar refresh token")
	}

	metadata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "falha ao criar sessão")
	}

	response := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return response, nil
}

func validateLoginUserRequest(request *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidatePassword(request.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
