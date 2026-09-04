package outbound

import (
	"context"

	"github.com/kidx45/Debter/internal/domain"
)

type UpdateContents interface {
	UpdateAccount(ctx context.Context, accountID, userID, amount int64) error
	InsertEntry(ctx context.Context, accountID, userID, amount int64, entryType, category string) (domain.Entry, error)
	UpdateBalanceTx(ctx context.Context, accountID, userID, amount int64, entryType, category string) (domain.Entry, domain.Account, error)
}
