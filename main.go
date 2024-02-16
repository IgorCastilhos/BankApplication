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
    "github.com/IgorCastilhos/BankApplication/worker"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/hibiken/asynq"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rakyll/statik/fs"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "google.golang.org/protobuf/encoding/protojson"
    "net"
    "net/http"
    "os"
)

func main() {
    config, err := utils.LoadConfig(".")
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível carregar as configurações")
    }
    
    if config.Environment == "development" {
        // Adiciona um log legível e colorido, usado principalmente em Development
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }
    
    connPool, err := pgxpool.New(context.Background(), config.DBSource)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível conectar ao banco de dados:")
    }
    
    runDBMigration(config.MigrationURL, config.DBSource)
    
    store := db.NewStore(connPool)
    
    redisOpt := asynq.RedisClientOpt{
        Addr: config.RedisAddress,
    }
    
    taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
    go runTaskProcessor(redisOpt, store)
    go runGatewayServer(config, store, taskDistributor)
    runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL, dbSource string) {
    migration, err := migrate.New(migrationURL, dbSource)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar uma nova migration:")
    }
    
    if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        log.Fatal().Err(err).Msg("falhou ao rodar migration up")
    }
    
    log.Info().Msg("migration realizada com sucesso")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
    taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
    log.Info().Msg("iniciando processador de tarefas")
    err := taskProcessor.Start()
    if err != nil {
        log.Fatal().Err(err).Msg("falha ao iniciar processador de tarefas")
    }
}

func runGrpcServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := grpcApi.NewServer(config, store, taskDistributor)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar o servidor:")
    }
    
    grpcLogger := grpc.UnaryInterceptor(grpcApi.GrpcLogger)
    grpcServer := grpc.NewServer(grpcLogger)
    pb.RegisterBankServer(grpcServer, server)
    reflection.Register(grpcServer)
    
    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar um listener:")
    }
    
    log.Info().Msgf("iniciando servidor gRPC na porta %s", listener.Addr().String())
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível inicializar servidor gRPC:")
    }
}

func runGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := grpcApi.NewServer(config, store, taskDistributor)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar o servidor:")
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
        log.Fatal().Err(err).Msg("não foi possível registar um server handler:")
    }
    
    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)
    
    statikFileServer, err := fs.New()
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar um sistema de arquivos:")
    }
    
    swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFileServer))
    mux.Handle("/swagger/", swaggerHandler)
    
    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar um listener:")
    }
    
    log.Info().Msgf("iniciando servidor HTTP Gateway na porta %s", listener.Addr().String())
    handler := grpcApi.HTTPLogger(mux)
    err = http.Serve(listener, handler)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível iniciar servidor HTTP Gateway:")
    }
}

func runGinServer(config utils.Config, store db.Store) {
    server, err := api.NewServer(config, store)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível criar o servidor:")
    }
    
    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("não foi possível conectar ao banco de dados:")
    }
}
