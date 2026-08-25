package outbound

import (
	"context"
	"time"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type postgresEntryRepo struct {
	queries *db.Queries
}

func NewPostgresEntryRepository(queries *db.Queries) outbound.EntryRepository {
	return &postgresEntryRepo{queries: queries}
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
