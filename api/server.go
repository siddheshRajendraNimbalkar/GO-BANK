package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
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
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
