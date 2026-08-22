package outbound

import (
	"context"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
)

type EntryRepository interface {
	GetEntriesByAccountId(ctx context.Context, arg db.GetEntriesByAccountIdParams) ([]db.Entry, error)
	FilterEntriesByDate(ctx context.Context, arg db.FilterEntriesByDateParams) ([]db.Entry, error)
	GetEntriesByCategoryAndType(ctx context.Context, arg db.GetEntriesByCategoryAndTypeParams) ([]db.GetEntriesByCategoryAndTypeRow, error)
}
