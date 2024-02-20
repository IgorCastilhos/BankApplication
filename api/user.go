package api

import (
    "errors"
    db "github.com/IgorCastilhos/BankApplication/db/sqlc"
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
    "time"
)

type createUserRequest struct {
    Username string `json:"username" binding:"required,alphanum"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
    Username          string    `json:"username"`
    FullName          string    `json:"full_name"`
    Email             string    `json:"email"`
    PasswordChangedAt time.Time `json:"password_changed_at"`
    CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
    return userResponse{
        Username:          user.Username,
        FullName:          user.FullName,
        Email:             user.Email,
        PasswordChangedAt: user.PasswordChangedAt,
        CreatedAt:         user.CreatedAt,
    }
}

func (server *Server) createUser(ctx *gin.Context) {
    var request createUserRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }
    
    hashedPassword, err := utils.HashPassword(request.Password)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    
    // Se for válido...
    arg := db.CreateUserParams{
        Username:       request.Username,
        HashedPassword: hashedPassword,
        FullName:       request.FullName,
        Email:          request.Email,
    }
    
    user, err := server.store.CreateUser(ctx, arg)
    if err != nil {
        if db.ErrorCode(err) == db.UniqueViolation {
            ctx.JSON(http.StatusForbidden, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    
    response := newUserResponse(user)
    // Se nenhum erro ocorrer, retornará OK com o usuário criado para o client
    ctx.JSON(http.StatusOK, response)
}

type loginUserRequest struct {
    Username string `json:"username" binding:"required,alphanum"`
    Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
    SessionID             uuid.UUID    `json:"session_id"`
    AccessToken           string       `json:"access_token,omitempty"`
    AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
    RefreshToken          string       `json:"refresh_token,omitempty"`
    RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at,omitempty"`
    User                  userResponse `json:"user,omitempty"`
}

func (server *Server) loginUser(ctx *gin.Context) {
    var req loginUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }
    
    user, err := server.store.GetUser(ctx, req.Username)
    if err != nil {
        if errors.Is(err, db.ErrRecordNotFound) {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    
    err = utils.CheckPassword(req.Password, user.HashedPassword)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }
    
    accessToken, accessPayload, err := server.tokenMaker.CreateToken(
        user.Username, user.Role, server.config.AccessTokenDuration,
    )
    
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    
    refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
        user.Username,
        user.Role,
        server.config.RefreshTokenDuration,
    )
    
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    
    session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
        ID:           refreshPayload.ID,
        Username:     user.Username,
        RefreshToken: refreshToken,
        UserAgent:    ctx.Request.UserAgent(),
        ClientIp:     ctx.ClientIP(),
        IsBlocked:    false,
        ExpiresAt:    refreshPayload.ExpiredAt,
    })
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    
    response := loginUserResponse{
        SessionID:             session.ID,
        AccessToken:           accessToken,
        AccessTokenExpiresAt:  accessPayload.ExpiredAt,
        RefreshToken:          refreshToken,
        RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
        User:                  newUserResponse(user),
    }
    ctx.JSON(http.StatusOK, response)
}
