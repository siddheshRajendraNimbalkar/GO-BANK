package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

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

	go func() {
		server, err := grpcapi.NewGRPCService(config, *store)
		if err != nil {
			log.Fatal("cannot create server GATEWAY", err)
		}
		grpcMux := runtime.NewServeMux()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

		if err != nil {
			log.Fatal("cannot create server GATEWAY", err)
		}
		mux := http.NewServeMux()
		mux.Handle("/", grpcMux)
		listener, err := net.Listen("tcp", config.Addr)

		if err != nil {
			log.Fatal("[Listener error GATEWAY]:", err)
		}
		log.Println("Listener GATEWAY", listener.Addr().String())

		err = http.Serve(listener, mux)
		if err != nil {
			log.Fatal("[error GATEWAY]:", err)
		}
	}()

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
