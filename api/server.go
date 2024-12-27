package api

import (
	"github.com/siddheshRajendraNimbalkar/GO-BANK/middleware"
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
	tokenMaker, err := token.NewPasetoMaker(config.Secret)
	if err != nil {
		panic(err) // Handle this better in production (e.g., log the error)
	}
	server.tokenMaker = tokenMaker
	router := gin.Default()

	//User
	router.POST("/user/sign-up", server.createUser)
	router.POST("/user/sign-in", server.compareUser)
	router.GET("/user/:userName", server.getUser)

	authRoute := router.Group("/").Use(middleware.AuthMiddleWare(tokenMaker))

	// Accounts
	authRoute.POST("/accounts", server.createAcount)
	authRoute.GET("/accounts/:id", server.GetAcount)
	authRoute.GET("/accounts", server.listAccounts)
	authRoute.GET("/accounts/delete/:id", server.DeleteAccounts)
	authRoute.GET("/accounts/update/:id", server.UpdateAccount)

	// Transfers
	authRoute.POST("transfer", server.Transfer)

	server.router = router
	server.config = config
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
