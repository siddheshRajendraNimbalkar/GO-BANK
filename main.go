package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
	grpcapi "github.com/siddheshRajendraNimbalkar/GO-BANK/gRPC_API"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/pb"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("[Config]::error in env", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("[ERROR IN Main_Test]::While connecting db", err)
	}

	store := db.NewStore(conn)

	// normal api

	// server := api.NewServer(config, store)

	// err = server.Start(config.Addr)

	// if err != nil {
	// 	log.Fatal("[ERROR OCCURE WHILE CONNECTING THE PORT]::", err.Error())
	// }

	// gRPC

	grpcService := grpc.NewServer()
	server, err := grpcapi.NewGRPCService(config, *store)

	if err != nil {
		log.Fatal("[main_grpc]:", err)
	}

	pb.RegisterSimpleBankServer(grpcService, server)

	reflection.Register(grpcService)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")

	if err != nil {
		log.Fatal("[Listener error]:", err)
	}

	log.Println("Listener ", listener.Addr().String())
	err = grpcService.Serve(listener)

	if err != nil {
		log.Fatal("[error]:", err)
	}
}
