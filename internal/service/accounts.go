package service

import (
	"context"
	"fmt"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type AccountService struct {
	AccountRepo outbound.AccountRepository
}

func NewAccountService(accountRepo outbound.AccountRepository) *AccountService {
	return &AccountService{
		AccountRepo: accountRepo,
	}
}

func (s *AccountService) GetAccountsByUserId(ctx context.Context, userID int64) ([]domain.Account, error) {
	res, err := s.AccountRepo.GetAccountsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no accounts found for user")
	}
	return res, nil
}
