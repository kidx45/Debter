package outbound

import (
	"database/sql"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
)

func NewPostgresUserRepository(conn *sql.DB) outbound.UserRepository {
	return &outbound.SqlRepository{
		Db: conn,
		Queries: db.New(conn),
	}
}