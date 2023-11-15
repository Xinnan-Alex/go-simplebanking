package api

import (
	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// endpoints for account
	router.POST("/accounts", server.createAccountAPI)
	router.GET("/accounts/:id", server.getAccountAPI)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccountAPI)

	// endpoints for transfer
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorReponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
