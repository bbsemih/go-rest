package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/bbsemih/gobank/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Can't establish connection to the Postgres:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
