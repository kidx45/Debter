package util

import (
	"context"
	"database/sql"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
)

func ExecTx(ctx context.Context, conn *sql.DB, fn func(queries *db.Queries) error) error {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.New(tx)
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
