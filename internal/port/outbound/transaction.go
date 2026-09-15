package outbound

import (
	"context"

	"github.com/kidx45/Debter/internal/domain"
)

type UpdateContents interface {
	UpdateBalanceTx(ctx context.Context, accountID, userID, amount int64, entryType, category string) (domain.Entry, domain.Account, error)
}
