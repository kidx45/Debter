package postgres_test

import (
	"context"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T, query *db.Queries, args ...db.CreateAccountParams) (db.Account, *db.Queries) {
	var arg db.CreateAccountParams

	if len(args) == 0 && query == nil {
		user, q := CreateRandomUser(t)
		arg = db.CreateAccountParams{
			UserID:        user.ID,
			AccountType:   "savings",
			AccountNumber: util.RandomNumber(100000, 999999),
			Balance:       5000,
		}
		query = q
	}

	account, err := query.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.UserID, account.UserID)
	require.Equal(t, "savings", account.AccountType)
	require.Equal(t, int64(5000), account.Balance)
	return account, query
}

func TestCreateAccount(t *testing.T, arg ...db.CreateAccountParams) {
	CreateRandomAccount(t, nil)
}

func TestGetAccountsByUserId(t *testing.T) {
	user, q := CreateRandomUser(t)

	createTestAccountDirect(t, user.ID, "savings", 10000)
	createTestAccountDirect(t, user.ID, "checking", 5000)
	createTestAccountDirect(t, user.ID, "investment", 15000)

	accounts, err := q.GetAccountsByUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, accounts, 3)

	require.Equal(t, "checking", accounts[0].AccountType)
	require.Equal(t, "investment", accounts[1].AccountType)
	require.Equal(t, "savings", accounts[2].AccountType)

	require.Equal(t, int64(5000), accounts[0].Balance)
	require.Equal(t, int64(15000), accounts[1].Balance)
	require.Equal(t, int64(10000), accounts[2].Balance)
}

func TestGetAccountsByUserIdNoAccounts(t *testing.T) {
	user := createTestUserDirect(t)

	q := TxQueriesTest(t)
	accounts, err := q.GetAccountsByUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.Empty(t, accounts)
}

func TestGetAccountsByUserIdBalance(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 7500)

	q := TxQueriesTest(t)
	accounts, err := q.GetAccountsByUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, account.ID, accounts[0].ID)
	require.Equal(t, int64(7500), accounts[0].Balance)

	err = q.CreditAccount(context.Background(), db.CreditAccountParams{
		Balance: 2500,
		ID:      account.ID,
		UserID:  user.ID,
	})
	require.NoError(t, err)

	accounts, err = q.GetAccountsByUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, int64(10000), accounts[0].Balance)
}
