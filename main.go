package main

import (
	"context"
	"github.com/IgorCastilhos/BankApplication/api"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/grpcApi"
	"github.com/IgorCastilhos/BankApplication/pb"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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
	runGrpcServer(config, store)
}

func runGrpcServer(config utils.Config, store db.Store) {
	server, err := grpcApi.NewServer(config, store)
	if err != nil {
		log.Fatal("não foi possível criar o servidor:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("não consegue criar um listener")
	}
	log.Printf("iniciando servidor gRPC em %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("não consegue iniciar servidor gRPC")
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("não foi possível criar o servidor:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("não foi possível conectar ao banco de dados:", err)
	}
}
