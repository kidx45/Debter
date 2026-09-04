package outbound

import (
	"context"
	"time"

	"github.com/kidx45/Debter/internal/domain"
)

type EntryRepository interface {
	CreateEntry(ctx context.Context, accountID, userID int64, amount int64, entryType, category string) (domain.Entry, error)
	Transfer(ctx context.Context, fromAccountID, toAccountID, userID int64, amount int64, entryType, category string) (domain.Entry, domain.Entry, error)
	GetEntriesByAccountId(ctx context.Context, accountID, userID int64) ([]domain.Entry, error)
	FilterEntriesByDate(ctx context.Context, accountID, userID int64, from, to time.Time) ([]domain.Entry, error)
	GetEntriesByCategoryAndType(ctx context.Context, accountID, userID int64, entryType string) ([]domain.CategorySummary, error)
}
