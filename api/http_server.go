package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/thehaung/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewHttpServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil
		}
	}

	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/accounts", server.createAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"errorMessage": err.Error()}
}
