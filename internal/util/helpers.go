package util

import (
	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/domain"
)

func DbEntryToDomain(e db.Entry) domain.Entry {
	return domain.Entry{
		ID:        e.ID,
		AccountID: e.AccountID,
		Amount:    e.Amount,
		Type:      e.Type,
		Category:  e.Category,
		CreatedAt: e.CreatedAt,
	}
}

func DbAccountToDomain(a db.Account) domain.Account {
	return domain.Account{
		ID:            a.ID,
		UserID:        a.UserID,
		AccountType:   a.AccountType,
		AccountNumber: a.AccountNumber,
		Balance:       a.Balance,
		CreatedAt:     a.CreatedAt,
	}
}
