package api

import (
	"fmt"

	"github.com/lfvm/simplebank/token"
	"github.com/lfvm/simplebank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/lfvm/simplebank/db/sqlc"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {

	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/login", server.logIn)
	router.POST("/refresh-token", server.refreshToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountById)
	authRoutes.GET("/accounts/", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)
	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
