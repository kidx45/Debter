package outbound

import (
	"context"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
)

type AccountRepository interface {
	GetAccountsByUserId(ctx context.Context, userID int64) ([]db.Account, error)
}
