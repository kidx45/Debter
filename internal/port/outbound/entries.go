package outbound

import (
	"context"
	"time"

	"github.com/kidx45/Debter/internal/domain"
)

type EntryRepository interface {
	GetEntriesByAccountId(ctx context.Context, accountID, userID int64) ([]domain.Entry, error)
	FilterEntriesByDate(ctx context.Context, accountID, userID int64, from, to time.Time) ([]domain.Entry, error)
	GetEntriesByCategoryAndType(ctx context.Context, accountID, userID int64, entryType string) ([]domain.CategorySummary, error)
}
