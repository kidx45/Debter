package postgres_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadEnv("../../test.env")
	if err != nil {
		log.Fatalf("Error loading the Env Configuration: %s", err)
	}

	testDB, err = sql.Open(config.DB_DRIVER_NAME, config.DB_URL)
	if err != nil {
		log.Fatalf("Error setting up db connection: %s", err)
	}

	exitCode := m.Run()
	testDB.Close()
	os.Exit(exitCode)
}

func TxQueriesTest (t *testing.T) *db.Queries {
	tx, err := testDB.BeginTx(context.Background(),nil)
	require.NoError(t,err)
	t.Cleanup(func() {
		tx.Rollback()
	})
	return db.New(tx)
}
