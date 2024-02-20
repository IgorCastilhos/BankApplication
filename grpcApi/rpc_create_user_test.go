package grpcApi

import (
    "context"
    "fmt"
    mockdb "github.com/IgorCastilhos/BankApplication/db/mock"
    db "github.com/IgorCastilhos/BankApplication/db/sqlc"
    "github.com/IgorCastilhos/BankApplication/pb"
    "github.com/IgorCastilhos/BankApplication/utils"
    mockwk "github.com/IgorCastilhos/BankApplication/worker/mock"
    "github.com/IgorCastilhos/toolkit/v2"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/require"
    "reflect"
    "testing"
)

var tools toolkit.Tools

type eqCreateUserTxParamsMatcher struct {
    arg      db.CreateUserTxParams
    password string
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
    actualArg, ok := x.(db.CreateUserTxParams)
    if !ok {
        return false
    }
    
    err := utils.CheckPassword(expected.password, actualArg.HashedPassword)
    if err != nil {
        return false
    }
    
    expected.arg.HashedPassword = actualArg.HashedPassword
    return reflect.DeepEqual(expected.arg, actualArg)
}

func (e eqCreateUserTxParamsMatcher) String() string {
    return fmt.Sprintf("arguemnto correspondente %v e senha %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string) gomock.Matcher {
    return eqCreateUserTxParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user db.User, password string) {
    password = tools.RandomString(6)
    hashedPassword, err := utils.HashPassword(password)
    require.NoError(t, err)
    
    user = db.User{
        Username:       utils.RandomOwner(),
        HashedPassword: hashedPassword,
        FullName:       utils.RandomOwner(),
        Email:          utils.RandomEmail(),
    }
    return
}

func TestCreateUserAPI(t *testing.T) {
    user, password := randomUser(t)
    
    testCases := []struct {
        name          string
        body          *pb.CreateUserRequest
        buildStubs    func(store *mockdb.MockStore)
        checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
    }{
        {
            name: "OK",
            body: &pb.CreateUserRequest{
                Username: user.Username,
                Password: password,
                FullName: user.FullName,
                Email:    user.Email,
            },
            buildStubs: func(store *mockdb.MockStore) {
                arg := db.CreateUserTxParams{
                    CreateUserParams: db.CreateUserParams{
                        Username: user.Username,
                        FullName: user.FullName,
                        Email:    user.Email,
                    },
                }
                store.EXPECT().
                    CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password)).
                    Times(1).
                    Return(db.CreateUserTxResult{User: user}, nil)
            },
            checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
                require.NoError(t, err)
                require.NotNil(t, res)
                createdUser := res.GetUser()
                require.Equal(t, user.Username, createdUser.Username)
                require.Equal(t, user.FullName, createdUser.FullName)
                require.Equal(t, user.Email, createdUser.Email)
            },
        },
    }
    
    for i := range testCases {
        tc := testCases[i]
        
        t.Run(tc.name, func(t *testing.T) {
            storeCtrl := gomock.NewController(t)
            defer storeCtrl.Finish()
            store := mockdb.NewMockStore(storeCtrl)
            
            taskCtrl := gomock.NewController(t)
            defer taskCtrl.Finish()
            taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)
            
            tc.buildStubs(store)
            server := newTestServer(t, store, taskDistributor)
            
            res, err := server.CreateUser(context.Background(), tc.req)
            tc.checkResponse(t, res, err)
        })
    }
}
