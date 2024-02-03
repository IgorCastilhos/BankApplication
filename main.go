package main

import (
	"database/sql"
	"github.com/IgorCastilhos/BankApplication/api"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/utils"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("não pôde carregar as configurações", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("não foi possível conectar ao banco de dados:", err)
	}
}
