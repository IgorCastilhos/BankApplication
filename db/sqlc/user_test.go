package db

import (
    "context"
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/IgorCastilhos/toolkit/v2"
    "github.com/jackc/pgx/v5/pgtype"
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

// TestCreateAccount testa se a função CreateAccount funciona conforme esperado.
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

func TestUpdateUserOnlyFullName(t *testing.T) {
    oldUser := createRandomUser(t)
    
    newFullName := utils.RandomOwner()
    
    updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
        Username: oldUser.Username,
        FullName: pgtype.Text{
            String: newFullName,
            Valid:  true,
        },
    })
    require.NoError(t, err)
    require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
    require.Equal(t, newFullName, updatedUser.FullName)
    require.Equal(t, oldUser.Email, updatedUser.Email)
    require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateEmailOnlyFullName(t *testing.T) {
    oldUser := createRandomUser(t)
    
    newEmail := utils.RandomEmail()
    updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
        Username: oldUser.Username,
        Email: pgtype.Text{
            String: newEmail,
            Valid:  true,
        },
    })
    
    require.NoError(t, err)
    require.NotEqual(t, oldUser.Email, updatedUser.Email)
    require.Equal(t, newEmail, updatedUser.Email)
    require.Equal(t, oldUser.FullName, updatedUser.FullName)
    require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
    oldUser := createRandomUser(t)
    
    newEmail := utils.RandomEmail()
    updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
        Username: oldUser.Username,
        Email: pgtype.Text{
            String: newEmail,
            Valid:  true,
        },
    })
    
    require.NoError(t, err)
    require.NotEqual(t, oldUser.Email, updatedUser.Email)
    require.Equal(t, newEmail, updatedUser.Email)
    require.Equal(t, oldUser.FullName, updatedUser.FullName)
    require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserPasswordOnly(t *testing.T) {
    oldUser := createRandomUser(t)
    
    newPassword := tool.RandomString(6)
    newHashedPassword, err := utils.HashPassword(newPassword)
    require.NoError(t, err)
    
    updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
        Username: oldUser.Username,
        HashedPassword: pgtype.Text{
            String: newHashedPassword,
            Valid:  true,
        },
    })
    
    require.NoError(t, err)
    require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
    require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
    require.Equal(t, oldUser.FullName, updatedUser.FullName)
    require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
    oldUser := createRandomUser(t)
    
    newFullName := utils.RandomOwner()
    newEmail := utils.RandomEmail()
    newPassword := tool.RandomString(6)
    newHashedPassword, err := utils.HashPassword(newPassword)
    require.NoError(t, err)
    
    updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
        Username: oldUser.Username,
        FullName: pgtype.Text{
            String: newFullName,
            Valid:  true,
        },
        Email: pgtype.Text{
            String: newEmail,
            Valid:  true,
        },
        HashedPassword: pgtype.Text{
            String: newHashedPassword,
            Valid:  true,
        },
    })
    
    require.NoError(t, err)
    require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
    require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
    
    require.NotEqual(t, oldUser.Email, updatedUser.Email)
    require.Equal(t, newEmail, updatedUser.Email)
    
    require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
    require.Equal(t, newFullName, updatedUser.FullName)
}
