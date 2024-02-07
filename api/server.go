package api

import (
	"fmt"
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/IgorCastilhos/BankApplication/token"
	"github.com/IgorCastilhos/BankApplication/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server define as requisições HTTP para o serviço bancário
type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine // Router para enviar cada requisição para API ao manipulador correto
}

// NewServer cria um novo servidor HTTP e configura o roteamento
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar um token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Adiciona rotas ao roteador
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start executa o servidor HTTP em um endereço específico, recebido por parâmetro
func (server *Server) Start(address string) error {
	return server.router.Run(address)
	// Todo: Adicionar lógica de graceful shutdown (desligamento normal) com sinais https://www.youtube.com/watch?v=vgGreJPn3q0
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
