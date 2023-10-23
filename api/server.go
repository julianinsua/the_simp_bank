package api

import (
	"github.com/gin-gonic/gin"
	"github.com/julianinsua/the_simp_bank.git/internal/database"
)

/* Struct to generate a server*/
type Server struct {
	store  database.Store
	router *gin.Engine
}

/* Create a new server struct, add routes andd return the server instance */
func NewServer(store database.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Routes go here
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.getAccountList)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
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
