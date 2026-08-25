package outbound

import (
	"context"

	"github.com/kidx45/Debter/internal/domain"
)

type AccountRepository interface {
	GetAccountsByUserId(ctx context.Context, userID int64) ([]domain.Account, error)
}
