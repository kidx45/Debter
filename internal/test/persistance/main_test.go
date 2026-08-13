package persistance

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/db"

	"github.com/kidx45/Debter/internal/util"
)

var testdb *db.Queries

func TestMain(m *testing.M) {
	config, err := util.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading the Env Configuration: %s", err)
	}

	conn, err := sql.Open(config.DB_DRIVER_NAME, config.DB_URL_TEST)
	if err != nil {
		log.Fatalf("Error setting up db connection: %s", err)
	}

	testdb = db.New(conn)
	exitcode := m.Run()
	conn.Close()
	os.Exit(exitcode)
}
