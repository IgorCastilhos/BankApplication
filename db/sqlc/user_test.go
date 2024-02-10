package db

import (
	"context"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/IgorCastilhos/toolkit/v2"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var tool toolkit.Tools

// createRandomAccount cria uma conta com valores aleatórios e retorna essa conta.
// Utiliza funções do pacote utils para gerar valores aleatórios para os campos da conta.
// Realiza várias verificações para assegurar que a conta foi criada corretamente.
func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(tool.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// Test_CreateAccount testa se a função CreateAccount funciona conforme esperado.
// Verifica se a conta criada corresponde aos parâmetros fornecidos e se os campos obrigatórios estão presentes.
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// Test_GetAccount testa a função GetAccount.
// Cria uma conta, recupera essa conta e verifica se os dados recuperados correspondem aos dados da conta criada.
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
