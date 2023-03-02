package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
)

const (
	dbDriver = "pgx"
	dbSource = "postgres://root:aregbesola@127.0.0.1:15432/omnibank?sslmode=disable"
)

var testQueries *Queries
var TestDB *sql.DB



func TestMain(m *testing.M) {
	var err error

	TestDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("Can't connect with database:", err)
	}

	testQueries = New(TestDB)

	os.Exit(m.Run())
}