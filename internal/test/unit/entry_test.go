package unit

import (
	"context"
	"testing"
	"time"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/service"
	mockrepository "github.com/kidx45/Debter/internal/test/mock/repository"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func RandomEntry(t *testing.T) db.Entry {
	return db.Entry{
		ID:        util.RandomNumber(1, 1000),
		AccountID: util.RandomNumber(1, 100),
		Amount:    util.RandomNumber(100, 10000),
		Type:      "expense",
		Category:  "food",
		CreatedAt: time.Now(),
	}
}

func NewTestEntryService(t *testing.T, DB outbound.EntryRepository) *service.EntryService {
	return service.NewEntryService(DB)
}

func TestGetEntriesByAccountId(t *testing.T) {
	entries := []db.Entry{RandomEntry(t), RandomEntry(t)}
	accountID := entries[0].AccountID

	testCases := []struct {
		name            string
		accountID       int64
		buildRepository func(repository *mockrepository.MockEntryRepository)
		checkResponse   func(t *testing.T, res []db.Entry, err error)
	}{
		{
			name:      "OK",
			accountID: accountID,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByAccountId(gomock.Any(), gomock.Eq(accountID)).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, res []db.Entry, err error) {
				require.NoError(t, err)
				require.Equal(t, entries, res)
			},
		},
		{
			name:      "Empty",
			accountID: accountID,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByAccountId(gomock.Any(), gomock.Eq(accountID)).Return([]db.Entry{}, nil)
			},
			checkResponse: func(t *testing.T, res []db.Entry, err error) {
				require.Error(t, err)
				require.Equal(t, "no entries found for this account ID", err.Error())
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockEntryRepository(ctrl)
			testCases[i].buildRepository(repository)
			EntryService := NewTestEntryService(t, repository)
			res, err := EntryService.GetEntriesByAccountId(context.Background(), testCases[i].accountID)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestFilterEntriesByDate(t *testing.T) {
	entries := []db.Entry{RandomEntry(t)}
	arg := db.FilterEntriesByDateParams{
		AccountID:   entries[0].AccountID,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		CreatedAt_2: time.Now(),
	}

	testCases := []struct {
		name            string
		arg             db.FilterEntriesByDateParams
		buildRepository func(repository *mockrepository.MockEntryRepository)
		checkResponse   func(t *testing.T, res []db.Entry, err error)
	}{
		{
			name: "OK",
			arg:  arg,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().FilterEntriesByDate(gomock.Any(), gomock.Eq(arg)).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, res []db.Entry, err error) {
				require.NoError(t, err)
				require.Equal(t, entries, res)
			},
		},
		{
			name: "Empty",
			arg:  arg,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().FilterEntriesByDate(gomock.Any(), gomock.Eq(arg)).Return([]db.Entry{}, nil)
			},
			checkResponse: func(t *testing.T, res []db.Entry, err error) {
				require.NoError(t, err)
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockEntryRepository(ctrl)
			testCases[i].buildRepository(repository)
			EntryService := NewTestEntryService(t, repository)
			res, err := EntryService.FilterEntriesByDate(context.Background(), testCases[i].arg)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestGetEntriesByCategoryAndType(t *testing.T) {
	results := []db.GetEntriesByCategoryAndTypeRow{
		{Category: "food", Total: 5000},
		{Category: "transport", Total: 2000},
	}
	arg := db.GetEntriesByCategoryAndTypeParams{
		AccountID: util.RandomNumber(1, 100),
		Type:      "expense",
	}

	testCases := []struct {
		name            string
		arg             db.GetEntriesByCategoryAndTypeParams
		buildRepository func(repository *mockrepository.MockEntryRepository)
		checkResponse   func(t *testing.T, res []db.GetEntriesByCategoryAndTypeRow, err error)
	}{
		{
			name: "OK",
			arg:  arg,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByCategoryAndType(gomock.Any(), gomock.Eq(arg)).Return(results, nil)
			},
			checkResponse: func(t *testing.T, res []db.GetEntriesByCategoryAndTypeRow, err error) {
				require.NoError(t, err)
				require.Equal(t, results, res)
			},
		},
		{
			name: "Empty",
			arg:  arg,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByCategoryAndType(gomock.Any(), gomock.Eq(arg)).Return([]db.GetEntriesByCategoryAndTypeRow{}, nil)
			},
			checkResponse: func(t *testing.T, res []db.GetEntriesByCategoryAndTypeRow, err error) {
				require.Error(t, err)
				require.Equal(t, "no entries found for this account ID", err.Error())
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockEntryRepository(ctrl)
			testCases[i].buildRepository(repository)
			EntryService := NewTestEntryService(t, repository)
			res, err := EntryService.GetEntriesByCategoryAndType(context.Background(), testCases[i].arg)
			testCases[i].checkResponse(t, res, err)
		})
	}
}
