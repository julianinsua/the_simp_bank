package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/julianinsua/the_simp_bank/internal/database"
	"github.com/julianinsua/the_simp_bank/token"
	"github.com/julianinsua/the_simp_bank/util"
	"github.com/pkg/errors"
)

/* Struct to generate a server*/
type Server struct {
	store      database.Store
	tokenMaker token.PASETOMaker
	router     *gin.Engine
	config     util.Config
}

/* Create a new server struct, add routes andd return the server instance */
func NewServer(config util.Config, store database.Store) (*Server, error) {
	tokenMaker, err := token.NewPASETOMaker(config.SymetricKey)
	if err != nil {
		return nil, errors.Errorf("couldn't initialize JWT token generator: %v", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	// Custom validation bindings
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

/*
Includes all the handlers on their specific routes
*/
func (srv *Server) setupRouter() {
	router := gin.Default()

	// Routes go here
	router.POST("/accounts", srv.createAccount)
	router.POST("/users/login", srv.loginUser)

	// Authorized routes
	authRoutes := router.Group("/").Use(authMiddleware(srv.tokenMaker))

	authRoutes.GET("/accounts", srv.getAccountList)
	authRoutes.GET("/accounts/:id", srv.getAccount)
	authRoutes.POST("/transfers", srv.createTransfer)
	authRoutes.POST("/users", srv.createUser)

	srv.router = router
}

/*
Starts the http server on a specific address
*/
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

/*
Marshals an error into something gin can return to the user
returns gin.H: map[string]interface{}
*/
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
