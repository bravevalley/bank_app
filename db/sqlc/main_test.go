package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/dassyareg/bank_app/utils"
	_ "github.com/jackc/pgx/stdlib"
)

var testQueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load config file ", err)
	}

	TestDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("Can't connect with database:", err)
	}

	testQueries = New(TestDB)

	os.Exit(m.Run())
}
