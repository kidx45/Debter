package outbound

import (
	"context"
	"database/sql"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/util"
)

type postgresTransactionRepo struct {
	queries *db.Queries
	conn    *sql.DB
}

func NewPostgresTransactionRepository(queries *db.Queries, conn *sql.DB) outbound.UpdateContents {
	return &postgresTransactionRepo{queries: queries, conn: conn}
}

func (r *postgresTransactionRepo) UpdateAccount(ctx context.Context, accountID, userID, amount int64) error {
	return r.updateAccount(ctx, r.queries, accountID, userID, amount)
}

func (r *postgresTransactionRepo) InsertEntry(ctx context.Context, accountID, userID, amount int64, entryType, category string) (domain.Entry, error) {
	return r.insertEntry(ctx, r.queries, accountID, userID, amount, entryType, category)
}

func (r *postgresTransactionRepo) UpdateBalanceTx(ctx context.Context, accountID, userID, amount int64, entryType, category string) (domain.Entry, domain.Account, error) {
	var entry domain.Entry
	var account domain.Account

	err := util.ExecTx(ctx, r.conn, func(q *db.Queries) error {
		var txErr error

		entry, txErr = r.insertEntry(ctx, q, accountID, userID, amount, entryType, category)
		if txErr != nil {
			return txErr
		}

		txErr = r.updateAccount(ctx, q, accountID, userID, amount)
		if txErr != nil {
			return txErr
		}

		row, txErr := q.GetAccountByID(ctx, accountID)
		if txErr != nil {
			return txErr
		}
		account = dbAccountToDomain(row)

		return nil
	})
	if err != nil {
		return domain.Entry{}, domain.Account{}, err
	}

	return entry, account, nil
}

func (r *postgresTransactionRepo) updateAccount(ctx context.Context, q *db.Queries, accountID, userID, amount int64) error {
	if amount >= 0 {
		return q.CreditAccount(ctx, db.CreditAccountParams{
			Balance: amount,
			ID:      accountID,
			UserID:  userID,
		})
	}

	rows, err := q.DebitAccount(ctx, db.DebitAccountParams{
		Balance: amount,
		ID:      accountID,
		UserID:  userID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *postgresTransactionRepo) insertEntry(ctx context.Context, q *db.Queries, accountID, userID, amount int64, entryType, category string) (domain.Entry, error) {
	result, err := q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: accountID,
		Amount:    amount,
		Type:      entryType,
		Category:  category,
	})
	if err != nil {
		return domain.Entry{}, err
	}
	return dbEntryToDomain(result), nil
}