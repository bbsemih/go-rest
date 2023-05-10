package api

import (
	"fmt"

	db "github.com/bbsemih/gobank/db/sqlc"
	"github.com/bbsemih/gobank/token"
	"github.com/bbsemih/gobank/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
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

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	r := gin.Default()

	r.POST("/accounts", server.createAccount)
	r.POST("/accounts/login", server.loginUser)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts", server.listAccounts)
	r.DELETE("/accounts/:id", server.deleteAccount)
	//update account?
	//role based access control for delete and update???

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
