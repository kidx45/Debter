package outbound

import (
	"context"
	"database/sql"
	"time"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/util"
)

type postgresEntryRepo struct {
	queries *db.Queries
	conn    *sql.DB
}

func NewPostgresEntryRepository(queries *db.Queries, conn *sql.DB) outbound.EntryRepository {
	return &postgresEntryRepo{queries: queries, conn: conn}
}

func (r *postgresEntryRepo) CreateEntry(ctx context.Context, accountID, userID int64, amount int64, entryType, category string) (domain.Entry, error) {
	result, err := r.queries.CreateEntry(ctx, db.CreateEntryParams{
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

func (r *postgresEntryRepo) Transfer(ctx context.Context, fromAccountID, toAccountID, userID int64, amount int64, entryType, category string) (domain.Entry, domain.Entry, error) {
	var fromEntry, toEntry db.Entry

	err := util.ExecTx(ctx, r.conn, func(q *db.Queries) error {
		var txErr error

		fromEntry, txErr = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: fromAccountID,
			Amount:    -amount,
			Type:      entryType,
			Category:  category,
		})
		if txErr != nil {
			return txErr
		}

		toEntry, txErr = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: toAccountID,
			Amount:    amount,
			Type:      entryType,
			Category:  category,
		})
		if txErr != nil {
			return txErr
		}

		rows, txErr := q.DebitAccount(ctx, db.DebitAccountParams{
			Balance: -amount,
			ID:      fromAccountID,
			UserID:  userID,
		})
		if txErr != nil {
			return txErr
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		txErr = q.CreditAccount(ctx, db.CreditAccountParams{
			Balance: amount,
			ID:      toAccountID,
			UserID:  userID,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return domain.Entry{}, domain.Entry{}, err
	}

	return dbEntryToDomain(fromEntry), dbEntryToDomain(toEntry), nil
}

func (r *postgresEntryRepo) GetEntriesByAccountId(ctx context.Context, accountID, userID int64) ([]domain.Entry, error) {
	results, err := r.queries.GetEntriesByAccountId(ctx, db.GetEntriesByAccountIdParams{
		AccountID: accountID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	entries := make([]domain.Entry, len(results))
	for i, e := range results {
		entries[i] = dbEntryToDomain(e)
	}
	return entries, nil
}

func (r *postgresEntryRepo) FilterEntriesByDate(ctx context.Context, accountID, userID int64, from, to time.Time) ([]domain.Entry, error) {
	results, err := r.queries.FilterEntriesByDate(ctx, db.FilterEntriesByDateParams{
		AccountID:   accountID,
		UserID:      userID,
		CreatedAt:   from,
		CreatedAt_2: to,
	})
	if err != nil {
		return nil, err
	}
	entries := make([]domain.Entry, len(results))
	for i, e := range results {
		entries[i] = dbEntryToDomain(e)
	}
	return entries, nil
}

func (r *postgresEntryRepo) GetEntriesByCategoryAndType(ctx context.Context, accountID, userID int64, entryType string) ([]domain.CategorySummary, error) {
	results, err := r.queries.GetEntriesByCategoryAndType(ctx, db.GetEntriesByCategoryAndTypeParams{
		AccountID: accountID,
		UserID:    userID,
		Type:      entryType,
	})
	if err != nil {
		return nil, err
	}
	summaries := make([]domain.CategorySummary, len(results))
	for i, s := range results {
		summaries[i] = domain.CategorySummary{
			Category: s.Category,
			Total:    s.Total,
		}
	}
	return summaries, nil
}

func dbEntryToDomain(e db.Entry) domain.Entry {
	return domain.Entry{
		ID:        e.ID,
		AccountID: e.AccountID,
		Amount:    e.Amount,
		Type:      e.Type,
		Category:  e.Category,
		CreatedAt: e.CreatedAt,
	}
}
