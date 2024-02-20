package grpcApi

import (
    db "github.com/IgorCastilhos/BankApplication/db/sqlc"
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/IgorCastilhos/BankApplication/worker"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
    config := utils.Config{
        TokenSymmetricKey:   tools.RandomString(32),
        AccessTokenDuration: time.Minute,
    }
    
    server, err := NewServer(config, store, taskDistributor)
    require.NoError(t, err)
    
    return server
}
