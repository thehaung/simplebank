package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/thehaung/simplebank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewHttpServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()
	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"errorMessage": err.Error()}
}
