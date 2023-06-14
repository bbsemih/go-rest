package api

import (
	"fmt"

	db "github.com/bbsemih/gobank/internal/db/sqlc"
	"github.com/bbsemih/gobank/pkg/token"
	"github.com/bbsemih/gobank/pkg/util"
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

	r.POST("/users", server.createUser)
	r.POST("/users/login", server.loginUser)

	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts", server.listAccounts)
	r.DELETE("/accounts/:id", server.deleteAccount)
	//TODO
	//update account???
	//role based access control for delete and update???

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
