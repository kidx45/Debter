package outbound

import (
	"database/sql"

	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/port/outbound/db"
)

func NewRepository(conn *sql.DB) outbound.Repository {
	return &outbound.SqlRepository{
		Db: conn,
		Queries: db.New(conn),
	}
}