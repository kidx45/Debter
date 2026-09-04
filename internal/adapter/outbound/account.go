package outbound

import (
	"context"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type postgresAccountRepo struct {
	queries *db.Queries
}

func NewPostgresAccountRepository(queries *db.Queries) outbound.AccountRepository {
	return &postgresAccountRepo{queries: queries}
}

func (r *postgresAccountRepo) GetAccountsByUserId(ctx context.Context, userID int64) ([]domain.Account, error) {
	results, err := r.queries.GetAccountsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	accounts := make([]domain.Account, len(results))
	for i, a := range results {
		accounts[i] = dbAccountToDomain(a)
	}
	return accounts, nil
}

func dbAccountToDomain(a db.Account) domain.Account {
	return domain.Account{
		ID:            a.ID,
		UserID:        a.UserID,
		AccountType:   a.AccountType,
		AccountNumber: a.AccountNumber,
		Balance:       a.Balance,
		CreatedAt:     a.CreatedAt,
	}
}
