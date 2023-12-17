package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/julianinsua/the_simp_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries = nil
var testDB *sql.DB = nil

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
