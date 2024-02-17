package db

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
    Querier
    TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
    CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// SQLStore fornece funções para executar consultas e transações SQL no banco de dados real
type SQLStore struct {
    connPool *pgxpool.Pool
    *Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
    return &SQLStore{
        connPool: connPool,
        Queries:  New(connPool),
    }
}
