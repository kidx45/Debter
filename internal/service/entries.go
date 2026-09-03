package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type EntryService struct {
	EntryRepo outbound.EntryRepository
}

func NewEntryService(entryRepo outbound.EntryRepository) *EntryService {
	return &EntryService{
		EntryRepo: entryRepo,
	}
}

func (s *EntryService) Transfer(ctx context.Context, fromAccountID, toAccountID, userID int64, amount int64, entryType, category string) (domain.Entry, domain.Entry, error) {
	if amount <= 0 {
		return domain.Entry{}, domain.Entry{}, fmt.Errorf("transfer amount must be positive")
	}
	return s.EntryRepo.Transfer(ctx, fromAccountID, toAccountID, userID, amount, entryType, category)
}

func (s *EntryService) GetEntriesByAccountId(ctx context.Context, accountID, userID int64) ([]domain.Entry, error) {
	res, err := s.EntryRepo.GetEntriesByAccountId(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no entries found for this account ID")
	}
	return res, nil
}

func (s *EntryService) FilterEntriesByDate(ctx context.Context, accountID, userID int64, from, to time.Time) ([]domain.Entry, error) {
	return s.EntryRepo.FilterEntriesByDate(ctx, accountID, userID, from, to)
}

func (s *EntryService) GetEntriesByCategoryAndType(ctx context.Context, accountID, userID int64, entryType string) ([]domain.CategorySummary, error) {
	res, err := s.EntryRepo.GetEntriesByCategoryAndType(ctx, accountID, userID, entryType)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no entries found for this account ID")
	}
	return res, nil
}
