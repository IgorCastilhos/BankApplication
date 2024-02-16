package grpcApi

import (
    "context"
    "errors"
    db "github.com/IgorCastilhos/BankApplication/db/sqlc"
    "github.com/IgorCastilhos/BankApplication/pb"
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/IgorCastilhos/BankApplication/validation"
    "github.com/jackc/pgx/v5/pgtype"
    "google.golang.org/genproto/googleapis/rpc/errdetails"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "time"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    authPayload, err := server.authorizeUser(ctx)
    if err != nil {
        return nil, unauthenticatedError(err)
    }
    
    violations := validateUpdateUserRequest(req)
    if violations != nil {
        return nil, invalidArgumentError(violations)
    }
    
    if authPayload.Username != req.GetUsername() {
        return nil, status.Errorf(codes.PermissionDenied, "não é possível atualizar informações de outros usuários")
    }
    
    // Se for válido...
    arg := db.UpdateUserParams{
        Username: req.GetUsername(),
        FullName: pgtype.Text{
            String: req.GetFullName(),
            Valid:  req.FullName != nil,
        },
        Email: pgtype.Text{
            String: req.GetEmail(),
            Valid:  req.Email != nil,
        },
    }
    
    if req.Password != nil {
        hashedPassword, err := utils.HashPassword(req.GetPassword())
        if err != nil {
            return nil, status.Errorf(codes.Internal, "falhou ao fazer o hash da senha: %s", err)
        }
        
        arg.HashedPassword = pgtype.Text{
            String: hashedPassword,
            Valid:  true,
        }
        
        arg.PasswordChangedAt = pgtype.Timestamptz{
            Time:  time.Now(),
            Valid: true,
        }
    }
    
    user, err := server.store.UpdateUser(ctx, arg)
    if err != nil {
        if errors.Is(err, db.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "usuário não encontrado")
        }
        return nil, status.Errorf(codes.Internal, "falhou ao criar um usuário: %s", err)
    }
    
    response := &pb.UpdateUserResponse{
        User: convertUser(user),
    }
    return response, nil
}

func validateUpdateUserRequest(request *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
    if err := validation.ValidateUsername(request.GetUsername()); err != nil {
        violations = append(violations, fieldViolation("username", err))
    }
    
    if request.Password != nil {
        if err := validation.ValidatePassword(request.GetPassword()); err != nil {
            violations = append(violations, fieldViolation("password", err))
        }
    }
    
    if request.FullName != nil {
        if err := validation.ValidateFullName(request.GetFullName()); err != nil {
            violations = append(violations, fieldViolation("full_name", err))
        }
    }
    
    if request.Email != nil {
        if err := validation.ValidateEmail(request.GetEmail()); err != nil {
            violations = append(violations, fieldViolation("email", err))
        }
    }
    
    return violations
}
