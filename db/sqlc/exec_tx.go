package db

import (
    "context"
    "fmt"
)

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
    tx, err := store.connPool.Begin(ctx)
    if err != nil {
        return err
    }
    
    q := New(tx)
    err = fn(q)
    if err != nil {
        if rollbackError := tx.Rollback(ctx); rollbackError != nil {
            return fmt.Errorf("transaction error: %v, rollback error: %v", err, rollbackError)
        }
        return err
    }
    return tx.Commit(ctx)
}
