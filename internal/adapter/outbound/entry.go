package outbound

import (
	"database/sql"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
)

func NewPostgresEntryRepository(conn *sql.DB) outbound.EntryRepository {
	return &outbound.SqlRepository{
		Db:      conn,
		Queries: db.New(conn),
	}
}
