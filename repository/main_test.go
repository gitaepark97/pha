package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gitaepark/pha/util"

	_ "github.com/go-sql-driver/mysql"
)

var (
	testQueries *Queries
	testDB      *sql.DB

	testConfig = util.Config{
		JWTSecret:            util.CreateRandomString(32),
		RefreshTokenDuration: time.Minute,
	}
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("./..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
