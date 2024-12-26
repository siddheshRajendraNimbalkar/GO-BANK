package api

import (
	"github.com/siddheshRajendraNimbalkar/GO-BANK/token"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"

	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

type Server struct {
	config     util.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Accounts
	router.POST("/accounts", server.createAcount)
	router.GET("/accounts/:id", server.GetAcount)
	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/delete/:id", server.DeleteAccounts)
	router.GET("/accounts/update/:id", server.UpdateAccount)

	// Transfers
	router.POST("transfer", server.Transfer)

	//User
	router.POST("/user/sign-up", server.createUser)
	router.POST("/user/sign-in", server.compareUser)
	router.GET("/user/:userName", server.getUser)

	server.router = router
	server.config = config
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
