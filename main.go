package main

import (
	"database/sql"
	"fmt"
	"log"

	DB "github.com/yankycranky/my-bank/db/sqlc"
)

var query *DB.Queries

const (
	driver = "postgres"
	source = "postgres://postgres:secret@localhost:5432/my_bank?sslmode=disable"
)

func main() {

	fmt.Println("This is main")
	dbtx, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal("Error connecting Database", err.Error())
	}
	query = DB.New(dbtx)
}

func GetQuery() *DB.Queries {
	return query
}
