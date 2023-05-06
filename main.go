package main

import (
	"database/sql"
	"log"

	"github.com/bbsemih/gobank/api"
	db "github.com/bbsemih/gobank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:secret@localhost:5432/gobank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cant't establish connection to the Postgres: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("Can't start the server: ", err)
	}
}
