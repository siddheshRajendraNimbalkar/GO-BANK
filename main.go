package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/siddheshRajendraNimbalkar/GO-BANK/api"
	db "github.com/siddheshRajendraNimbalkar/GO-BANK/db/sqlc"
)

const (
	dbDriv   = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
	addr     = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriv, dbSource)

	if err != nil {
		log.Fatal("[ERROR IN Main_Test]::While connecting db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(addr)

	if err != nil {
		log.Fatal("[ERROR OCCURE WHILE CONNECTING THE PORT]::", err.Error())
	}
}
