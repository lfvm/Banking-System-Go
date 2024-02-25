package db

import (
	"database/sql"
	"log"
	"testing"

	"github.com/lfvm/simplebank/utils"
	_ "github.com/lib/pq"
)



var testQueries *Queries
var testDb *sql.DB
var config utils.Config

func TestMain(m *testing.M){

	var err error

	config, err = utils.LoadConfig("../../")


	if err != nil {
		log.Fatal("Could not load env vars: ",err)
	}


	testDb,err = sql.Open(config.DbDriver,config.DbSource)

	if err != nil {
		log.Fatal("Could not connect to the database: ",err)
	}

	testQueries = New(testDb)

	m.Run()
}