package outbound

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
)

type Repository interface {
	db.Querier
}

type SqlRepository struct {
	*db.Queries
	Db *sql.DB
}

func (r *SqlRepository) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rollBackerror := tx.Rollback(); rollBackerror != nil {
			return fmt.Errorf("transaction err: %s, query err: %s", rollBackerror, err)
		}
		return err
	}
	return tx.Commit()
}