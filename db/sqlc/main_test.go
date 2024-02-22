package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M){

	var err error

	testDb,err = sql.Open(dbDriver,dbSource)

	if err != nil {
		log.Fatal("Could not connect to the database: ",err)
	}

	testQueries = New(testDb)

	m.Run()
}