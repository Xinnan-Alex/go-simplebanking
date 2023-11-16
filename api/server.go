package api

import (
	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// endpoints for account
	router.POST("/accounts", server.createAccountAPI)
	router.GET("/accounts/:id", server.getAccountAPI)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccountAPI)

	// endpoints for transfer
	router.POST("/transfers", server.createTransfer)

	//endpoint for users
	router.POST("/users", server.createUserAPI)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorReponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
