package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/yankycranky/my-bank/util"
)

const (
	driver = "postgres"
	source = "postgres://postgres:secret@localhost:5432/my_bank?sslmode=disable"
)

var testQueries *Queries

var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	dbConn, err = sql.Open(driver, source)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	util.Init()
	testQueries = New(dbConn)
	os.Exit(m.Run())
}
