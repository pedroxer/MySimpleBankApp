package api

import (
	db "MelBank/db/sqlc"
	_ "MelBank/docs"
	"MelBank/token"
	"MelBank/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server serves http requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//routes to router
	server.setupRouter()
	return server, nil

}
func (server *Server) setupRouter() {
	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	router.PUT("/accounts/update", server.updateAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/manager", server.createManager)
	router.POST("/manager/login", server.loginManager)
	authRoutes.GET("/manager/requests", server.listAllRequests)
	authRoutes.POST("/manager/request/", server.checkRequest)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router = router
}

// Start run the http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
