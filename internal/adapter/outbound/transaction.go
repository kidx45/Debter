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

func (r *postgresTransactionRepo) UpdateBalanceTx(ctx context.Context, accountID, userID, amount int64, entryType, category string) (domain.Entry, domain.Account, error) {
	var entry domain.Entry
	var account domain.Account

	err := util.ExecTx(ctx, r.conn, func(q *db.Queries) error {
		var txErr error

		entry, txErr = r.insertEntry(ctx, q, accountID, userID, amount, entryType, category)
		if txErr != nil {
			return txErr
		}

		account, txErr = r.updateAccount(ctx, q, accountID, userID, amount)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return domain.Entry{}, domain.Account{}, err
	}

	return entry, account, nil
}

func (r *postgresTransactionRepo) updateAccount(ctx context.Context, q *db.Queries, accountID, userID, amount int64) (domain.Account, error) {
	account, err := q.UpdateUserAccount(ctx, db.UpdateUserAccountParams{
		Balance: amount,
		ID:      accountID,
		UserID:  userID,
	})
	if err != nil {
		return domain.Account{}, err
	}
	return util.DbAccountToDomain(account), nil
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
	return util.DbEntryToDomain(result), nil
}
