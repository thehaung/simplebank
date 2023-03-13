package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thehaung/simplebank/config"
	db "github.com/thehaung/simplebank/db/sqlc"
	"github.com/thehaung/simplebank/token"
)

type Server struct {
	cfg        *config.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewHttpServer(cfg *config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(cfg.TokenSymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		cfg:        cfg,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}
	server.registerRouter()

	return server, nil
}

func (s *Server) registerRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	router.GET("/accounts", s.listAccount)
	router.GET("/accounts/:id", s.getAccount)
	router.POST("/accounts", s.createAccount)

	router.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"errorMessage": err.Error()}
}
