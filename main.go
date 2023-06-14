package main

import (
	"database/sql"
	"log"

	"github.com/bbsemih/gobank/api"
	db "github.com/bbsemih/gobank/internal/db/sqlc"
	"github.com/bbsemih/gobank/pkg/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cant't establish connection to the Postgres: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't start the server: ", err)
	}
}
