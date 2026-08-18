package service

import (
	"context"
	"fmt"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type AccountService struct {
	DB outbound.AccountRepository
}

func NewAccountService(DB outbound.AccountRepository) *AccountService {
	return &AccountService{
		DB: DB,
	}
}

func (s *AccountService) GetAccountsByUserId(ctx context.Context, userID int64) ([]db.Account, error) {
	res, err := s.DB.GetAccountsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no accounts found for user")
	}
	return res, nil
}
