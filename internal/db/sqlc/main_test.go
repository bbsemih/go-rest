package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/bbsemih/gobank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't establish connection to the Postgres:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
