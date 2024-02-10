package main

import (
	"context"
	"github.com/IgorCastilhos/BankApplication/api"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("não pôde carregar as configurações", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("não foi possível conectar ao banco de dados:", err)
	}

	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("não foi possível criar o servidor:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("não foi possível conectar ao banco de dados:", err)
	}
}
