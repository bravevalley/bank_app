package main

import (
	"database/sql"
	"log"

	"github.com/dassyareg/bank_app/api"
	db "github.com/dassyareg/bank_app/db/sqlc"
	_ "github.com/jackc/pgx/stdlib"
)

const (
	dbDriver = "pgx"
	dbSource = "postgres://root:aregbesola@127.0.0.1:15432/omnibank?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Could not connect database!")
	}

	masterQ := db.NewMasterQuery(conn)

	server := api.NewServer(masterQ)

	err = server.Start(address)
	if err != nil {
		log.Fatal("Cant start server")
	}

}
