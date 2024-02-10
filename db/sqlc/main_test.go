package db

import (
	"context"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("não pôde carregar as configurações", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("não pôde conectar ao banco de dados:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
