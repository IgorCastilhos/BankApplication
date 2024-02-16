package grpcApi

import (
    "context"
    "fmt"
    "github.com/IgorCastilhos/BankApplication/token"
    "google.golang.org/grpc/metadata"
    "strings"
)

const (
    authorizationHeader = "authorization"
    authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, fmt.Errorf("metadados não encontrados")
    }
    
    values := md.Get(authorizationHeader)
    if len(values) == 0 {
        return nil, fmt.Errorf("cabeçalho de autorização não encontrado")
    }
    
    // <authorization-type> <authorization-data>
    authHeader := values[0]
    fields := strings.Fields(authHeader)
    if len(fields) < 2 {
        return nil, fmt.Errorf("formato do cabeçalho de autorização inválido")
    }
    
    authType := strings.ToLower(fields[0])
    if authType != authorizationBearer {
        return nil, fmt.Errorf("tipo de autorização não suportado: %s", authType)
    }
    
    accessToken := fields[1]
    payload, err := server.tokenMaker.VerifyToken(accessToken)
    if err != nil {
        return nil, fmt.Errorf("token de acesso inválido: %s", err)
    }
    return payload, nil
}
