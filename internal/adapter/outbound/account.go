package outbound

import (
	"database/sql"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
)

func NewPostgresAccountRepository(conn *sql.DB) outbound.AccountRepository {
	return &outbound.SqlRepository{
		Db:      conn,
		Queries: db.New(conn),
	}
}
