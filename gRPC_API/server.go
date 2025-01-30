package grpcapi

import (
	"fmt"

	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/pb"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/token"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
)

type GRPCService struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewGRPCService(config util.Config, store db.Store) (*GRPCService, error) {
	tokenMaker, err := token.NewPasetoMaker(config.Secret)

	if err != nil {
		return nil, fmt.Errorf("[Token_Error] %w", err)
	}
	server := &GRPCService{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
