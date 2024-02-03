package api

import (
	db "github.com/IgorCastilhos/BankApplication/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server define as requisições HTTP para o serviço bancário
type Server struct {
	store  db.Store
	router *gin.Engine // Router para enviar cada requisição para API ao manipulador correto
}

// NewServer cria um novo servidor HTTP e configura o roteamento
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Adiciona rotas ao roteador
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start executa o servidor HTTP em um endereço específico, recebido por parâmetro
func (server *Server) Start(address string) error {
	return server.router.Run(address)
	// Todo: Adicionar lógica de graceful shutdown (desligamento normal) com sinais https://www.youtube.com/watch?v=vgGreJPn3q0
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
