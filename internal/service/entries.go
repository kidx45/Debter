package service

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
)

type EntryService struct {
	DB outbound.EntryRepository
}

func NewEntryService(DB outbound.EntryRepository) *EntryService {
	return &EntryService{
		DB: DB,
	}
}

func (s *EntryService) GetEntriesByAccountId(ctx context.Context, arg db.GetEntriesByAccountIdParams) ([]db.Entry, error) {
	res, err := s.DB.GetEntriesByAccountId(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no entries found for this account ID")
		}
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no entries found for this account ID")
	}
	return res, nil
}

func (s *EntryService) FilterEntriesByDate(ctx context.Context, arg db.FilterEntriesByDateParams) ([]db.Entry, error) {
	res, err := s.DB.FilterEntriesByDate(ctx, arg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *EntryService) GetEntriesByCategoryAndType(ctx context.Context, arg db.GetEntriesByCategoryAndTypeParams) ([]db.GetEntriesByCategoryAndTypeRow, error) {
	res, err := s.DB.GetEntriesByCategoryAndType(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no entries found for this account ID")
		}
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no entries found for this account ID")
	}
	return res, nil
}
