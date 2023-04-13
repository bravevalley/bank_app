package main

import (
	"database/sql"
	"log"

	"github.com/dassyareg/bank_app/api"
	db "github.com/dassyareg/bank_app/db/sqlc"
	"github.com/dassyareg/bank_app/utils"
	_ "github.com/lib/pq"
)

var err error

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config file ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Could not connect database!")
	}

	masterQ := db.NewMasterQuery(conn)

	server := api.NewServer(masterQ)

	err = server.Start(config.Address)
	if err != nil {
		log.Fatal("Cant start server")
	}

}
