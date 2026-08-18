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

func RandomAccount(t *testing.T) db.Account {
	return db.Account{
		ID:            util.RandomNumber(1, 1000),
		UserID:        util.RandomNumber(1, 100),
		AccountType:   "savings",
		AccountNumber: util.RandomNumber(100000, 999999),
		CreatedAt:     time.Now(),
	}
}

func NewTestAccountService(t *testing.T, DB outbound.AccountRepository) *service.AccountService {
	return service.NewAccountService(DB)
}

func TestGetAccountsByUserId(t *testing.T) {
	accounts := []db.Account{RandomAccount(t), RandomAccount(t)}
	userID := accounts[0].UserID

	testCases := []struct {
		name            string
		userID          int64
		buildRepository func(repository *mockrepository.MockAccountRepository)
		checkResponse   func(t *testing.T, res []db.Account, err error)
	}{
		{
			name:   "OK",
			userID: userID,
			buildRepository: func(repository *mockrepository.MockAccountRepository) {
				repository.EXPECT().GetAccountsByUserId(gomock.Any(), gomock.Eq(userID)).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, res []db.Account, err error) {
				require.NoError(t, err)
				require.Equal(t, accounts, res)
			},
		},
		{
			name:   "NotFound",
			userID: userID,
			buildRepository: func(repository *mockrepository.MockAccountRepository) {
				repository.EXPECT().GetAccountsByUserId(gomock.Any(), gomock.Eq(userID)).Return([]db.Account{}, nil)
			},
			checkResponse: func(t *testing.T, res []db.Account, err error) {
				require.Error(t, err)
				require.Equal(t, "no accounts found for user", err.Error())
				require.Empty(t, res)
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockAccountRepository(ctrl)
			testCases[i].buildRepository(repository)
			AccountService := NewTestAccountService(t, repository)
			res, err := AccountService.GetAccountsByUserId(context.Background(), testCases[i].userID)
			testCases[i].checkResponse(t, res, err)
		})
	}
}
