package api

import (
    "errors"
    "fmt"
    "github.com/IgorCastilhos/BankApplication/token"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

const (
    authorizationHeaderKey  = "authorization"
    authorizationTypeBearer = "bearer"
    authorizationPayloadKey = "authorization_payload"
)

// authMiddleware cria uma middleware (intermediário) do gin para autorização
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
        
        if len(authorizationHeader) == 0 {
            err := errors.New("cabeçalho de autorização não fornecido")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }
        
        fields := strings.Fields(authorizationHeader)
        if len(fields) < 2 {
            err := errors.New("formato de autorização de cabeçalho inválido")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }
        
        authorizationType := strings.ToLower(fields[0])
        if authorizationType != authorizationTypeBearer {
            err := fmt.Errorf("tipo de autorização não suportado %s", authorizationType)
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }
        
        accessToken := fields[1]
        payload, err := tokenMaker.VerifyToken(accessToken)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
            return
        }
        
        ctx.Set(authorizationPayloadKey, payload)
        ctx.Next()
    }
}
