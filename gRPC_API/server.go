package grpcapi

import (
	"fmt"

	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/token"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.Secret)

	if err != nil {
		return nil, fmt.Errorf("[Token_Error] %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
