package main

import (
    "context"
    "errors"
    "github.com/IgorCastilhos/BankApplication/api"
    db "github.com/IgorCastilhos/BankApplication/db/sqlc"
    _ "github.com/IgorCastilhos/BankApplication/doc/statik"
    "github.com/IgorCastilhos/BankApplication/grpcApi"
    "github.com/IgorCastilhos/BankApplication/pb"
    "github.com/IgorCastilhos/BankApplication/utils"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rakyll/statik/fs"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "google.golang.org/protobuf/encoding/protojson"
    "log"
    "net"
    "net/http"
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
    
    runDBMigration(config.MigrationURL, config.DBSource)
    
    store := db.NewStore(connPool)
    go runGatewayServer(config, store)
    runGrpcServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
    migration, err := migrate.New(migrationURL, dbSource)
    if err != nil {
        log.Fatal("não foi possível criar uma nova migration:", err)
    }
    
    if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        log.Fatal("falhou ao rodar migration up", err)
    }
    
    log.Println("migration realizada com sucesso")
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
        log.Fatal("não consegue criar um listener:", err)
    }
    log.Printf("iniciando servidor gRPC em %s", listener.Addr().String())
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatal("não consegue iniciar servidor gRPC:", err)
    }
}

func runGatewayServer(config utils.Config, store db.Store) {
    server, err := grpcApi.NewServer(config, store)
    if err != nil {
        log.Fatal("não foi possível criar o servidor:", err)
    }
    
    jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
        MarshalOptions: protojson.MarshalOptions{
            UseProtoNames: true,
        },
        UnmarshalOptions: protojson.UnmarshalOptions{
            DiscardUnknown: true,
        },
    })
    
    grpcMux := runtime.NewServeMux(jsonOption)
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)
    if err != nil {
        log.Fatal("Não foi possível registar um handler server:", err)
    }
    
    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)
    
    statikFileServer, err := fs.New()
    if err != nil {
        log.Fatal("Não pôde criar um sistema de arquivos:", err)
    }
    
    swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFileServer))
    mux.Handle("/swagger/", swaggerHandler)
    
    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal("não consegue criar um listener:", err)
    }
    log.Printf("iniciando servidor Gateway HTTP  em %s", listener.Addr().String())
    err = http.Serve(listener, mux)
    if err != nil {
        log.Fatal("não consegue iniciar servidor Gateway HTTP:", err)
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
