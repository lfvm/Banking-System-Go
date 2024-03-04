package main

import (
	"database/sql"
	"log"

	"github.com/lfvm/simplebank/api"
	db "github.com/lfvm/simplebank/db/sqlc"
	"github.com/lfvm/simplebank/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load env vars: ", err)
	}

	connection, err := sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	store := db.NewStore(connection)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not start server: ", err)
	}
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Can not start server: ", err)
	}

}
