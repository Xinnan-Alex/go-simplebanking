package api

import (
	"fmt"
	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/Xinnan-Alex/simplebank/token"
	"github.com/Xinnan-Alex/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//endpoint for users
	router.POST("/users", server.createUserAPI)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// endpoints for account
	authRoutes.POST("/accounts", server.createAccountAPI)
	authRoutes.GET("/accounts/:id", server.getAccountAPI)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccountAPI)

	// endpoints for transfer
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
