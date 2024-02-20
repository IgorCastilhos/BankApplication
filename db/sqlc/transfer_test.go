package db

import (
    "context"
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
    arg := CreateTransferParams{
        FromAccountID: account1.ID,
        ToAccountID:   account2.ID,
        Amount:        utils.RandomMoney(),
    }
    
    transfer, err := testStore.CreateTransfer(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, transfer)
    
    require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
    require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
    require.Equal(t, arg.Amount, transfer.Amount)
    
    require.NotZero(t, transfer.ID)
    require.NotZero(t, transfer.CreatedAt)
    require.NotZero(t, arg.Amount, transfer.Amount)
    
    return transfer
}

func TestCreateTransfer(t *testing.T) {
    account1 := createRandomAccount(t)
    account2 := createRandomAccount(t)
    createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
    account1 := createRandomAccount(t)
    account2 := createRandomAccount(t)
    transfer1 := createRandomTransfer(t, account1, account2)
    
    transfer2, err := testStore.GetTransfer(context.Background(), transfer1.ID)
    require.NoError(t, err)
    require.NotEmpty(t, transfer2)
    
    require.Equal(t, transfer1.ID, transfer2.ID)
    require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
    require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
    require.Equal(t, transfer1.Amount, transfer2.Amount)
    require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}
