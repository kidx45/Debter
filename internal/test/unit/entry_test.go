package unit

import (
	"context"
	"testing"
	"time"

	"github.com/kidx45/Debter/internal/domain"
	"github.com/kidx45/Debter/internal/service"
	mockrepository "github.com/kidx45/Debter/internal/test/mock/repository"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func RandomEntry(t *testing.T) domain.Entry {
	return domain.Entry{
		ID:        util.RandomNumber(1, 1000),
		AccountID: util.RandomNumber(1, 100),
		Amount:    util.RandomNumber(100, 10000),
		Type:      "expense",
		Category:  "food",
		CreatedAt: time.Now(),
	}
}

func NewTestEntryService(t *testing.T, entryRepo *mockrepository.MockEntryRepository) *service.EntryService {
	return service.NewEntryService(entryRepo)
}

func TestGetEntriesByAccountId(t *testing.T) {
	userID := util.RandomNumber(1, 100)
	entries := []domain.Entry{RandomEntry(t), RandomEntry(t)}
	accountID := entries[0].AccountID

	testCases := []struct {
		name            string
		accountID       int64
		userID          int64
		buildRepository func(repository *mockrepository.MockEntryRepository)
		checkResponse   func(t *testing.T, res []domain.Entry, err error)
	}{
		{
			name:      "OK",
			accountID: accountID,
			userID:    userID,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByAccountId(gomock.Any(), gomock.Eq(accountID), gomock.Eq(userID)).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, res []domain.Entry, err error) {
				require.NoError(t, err)
				require.Equal(t, entries, res)
			},
		},
		{
			name:      "Empty",
			accountID: accountID,
			userID:    userID,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByAccountId(gomock.Any(), gomock.Eq(accountID), gomock.Eq(userID)).Return([]domain.Entry{}, nil)
			},
			checkResponse: func(t *testing.T, res []domain.Entry, err error) {
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
			entryService := NewTestEntryService(t, repository)
			res, err := entryService.GetEntriesByAccountId(context.Background(), testCases[i].accountID, testCases[i].userID)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestFilterEntriesByDate(t *testing.T) {
	userID := util.RandomNumber(1, 100)
	entries := []domain.Entry{RandomEntry(t)}
	accountID := entries[0].AccountID
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	testCases := []struct {
		name            string
		accountID       int64
		userID          int64
		from            time.Time
		to              time.Time
		buildRepository func(repository *mockrepository.MockEntryRepository)
		checkResponse   func(t *testing.T, res []domain.Entry, err error)
	}{
		{
			name:      "OK",
			accountID: accountID,
			userID:    userID,
			from:      from,
			to:        to,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().FilterEntriesByDate(gomock.Any(), gomock.Eq(accountID), gomock.Eq(userID), gomock.Eq(from), gomock.Eq(to)).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, res []domain.Entry, err error) {
				require.NoError(t, err)
				require.Equal(t, entries, res)
			},
		},
		{
			name:      "Empty",
			accountID: accountID,
			userID:    userID,
			from:      from,
			to:        to,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().FilterEntriesByDate(gomock.Any(), gomock.Eq(accountID), gomock.Eq(userID), gomock.Eq(from), gomock.Eq(to)).Return([]domain.Entry{}, nil)
			},
			checkResponse: func(t *testing.T, res []domain.Entry, err error) {
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
			entryService := NewTestEntryService(t, repository)
			res, err := entryService.FilterEntriesByDate(context.Background(), testCases[i].accountID, testCases[i].userID, testCases[i].from, testCases[i].to)
			testCases[i].checkResponse(t, res, err)
		})
	}
}

func TestGetEntriesByCategoryAndType(t *testing.T) {
	userID := util.RandomNumber(1, 100)
	accountID := util.RandomNumber(1, 100)
	entryType := "expense"
	results := []domain.CategorySummary{
		{Category: "food", Total: 5000},
		{Category: "transport", Total: 2000},
	}

	testCases := []struct {
		name            string
		accountID       int64
		userID          int64
		entryType       string
		buildRepository func(repository *mockrepository.MockEntryRepository)
		checkResponse   func(t *testing.T, res []domain.CategorySummary, err error)
	}{
		{
			name:      "OK",
			accountID: accountID,
			userID:    userID,
			entryType: entryType,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByCategoryAndType(gomock.Any(), gomock.Eq(accountID), gomock.Eq(userID), gomock.Eq(entryType)).Return(results, nil)
			},
			checkResponse: func(t *testing.T, res []domain.CategorySummary, err error) {
				require.NoError(t, err)
				require.Equal(t, results, res)
			},
		},
		{
			name:      "Empty",
			accountID: accountID,
			userID:    userID,
			entryType: entryType,
			buildRepository: func(repository *mockrepository.MockEntryRepository) {
				repository.EXPECT().GetEntriesByCategoryAndType(gomock.Any(), gomock.Eq(accountID), gomock.Eq(userID), gomock.Eq(entryType)).Return([]domain.CategorySummary{}, nil)
			},
			checkResponse: func(t *testing.T, res []domain.CategorySummary, err error) {
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
			entryService := NewTestEntryService(t, repository)
			res, err := entryService.GetEntriesByCategoryAndType(context.Background(), testCases[i].accountID, testCases[i].userID, testCases[i].entryType)
			testCases[i].checkResponse(t, res, err)
		})
	}
}
