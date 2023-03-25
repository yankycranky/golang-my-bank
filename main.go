package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/yankycranky/my-bank/api"
	DB "github.com/yankycranky/my-bank/db/sqlc"
	db "github.com/yankycranky/my-bank/db/sqlc"
)

var query *DB.Queries

const (
	driver = "postgres"
	source = "postgres://postgres:secret@localhost:5432/my_bank?sslmode=disable"
)

func main() {

	fmt.Println("This is main")
	dbtx, err := sql.Open(driver, source)
	store := db.NewStore(dbtx)
	server := api.NewServer(store)
	if err != nil {
		log.Fatal("Error connecting Database", err.Error())
	}
	server.Start("0.0.0.0:3001")
	query = DB.New(dbtx)
}

func GetQuery() *DB.Queries {
	return query
}
