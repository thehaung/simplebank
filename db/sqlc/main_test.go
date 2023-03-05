package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	_dbDriver = "postgres"
	_dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var (
	_testQueries *Queries
	_testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	_testDB, err = sql.Open(_dbDriver, _dbSource)

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	_testQueries = New(_testDB)

	os.Exit(m.Run())
}
