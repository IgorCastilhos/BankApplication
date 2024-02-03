package db

import (
	"database/sql"
	"github.com/IgorCastilhos/BankApplication/utils"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("não pôde carregar as configurações", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("não pôde conectar ao banco de dados:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
