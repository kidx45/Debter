package outbound

import (
	"database/sql"
	"github.com/kidx45/Debter/internal/port/outbound/db"
)

type Repository interface {
	db.Querier
}

type SqlRepository struct {
	*db.Queries
	Db *sql.DB
}

