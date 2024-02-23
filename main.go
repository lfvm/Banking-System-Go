package main

import (
	"database/sql"
	"log"

	"github.com/lfvm/simplebank/api"
	db "github.com/lfvm/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)


func main(){

	connection,err := sql.Open(dbDriver,dbSource)

	if err != nil {
		log.Fatal("Could not connect to the database: ",err)
	}

	store:= db.NewStore(connection)
	server := api.NewServer(store)
	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Can not start server: ",err)
	}

}