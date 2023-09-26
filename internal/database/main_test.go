package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries = nil
var testDB *sql.DB = nil

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/simp_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("unable to connect to database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
