package postgres_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T, q *db.Queries) db.User {
	user, err := q.CreateUser(context.Background(), db.CreateUserParams{
		Username:       util.RandomString(6),
		FullName:       util.RandomString(6),
		Email:          util.RandomString(6),
		HashedPassword: util.RandomString(6),
	})
	require.NoError(t, err)
	return user
}

func createTestAccount(t *testing.T, q *db.Queries, userID int64, initialBalance int64) db.Account {
	account, err := q.CreateAccount(context.Background(), db.CreateAccountParams{
		UserID:        userID,
		AccountType:   "checking",
		AccountNumber: util.RandomNumber(100000, 999999),
		Balance:       initialBalance,
	})
	require.NoError(t, err)
	return account
}

func TestCreateEntry(t *testing.T) {
	q := TxQueriesTest(t)
	user := createTestUser(t, q)
	account := createTestAccount(t, q, user.ID, 0)

	entry, err := q.CreateEntry(context.Background(), db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    5000,
		Type:      "expense",
		Category:  "food",
	})
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, int64(5000), entry.Amount)
	require.Equal(t, "expense", entry.Type)
	require.Equal(t, "food", entry.Category)
}

func TestTransferConsistency(t *testing.T) {
	q := TxQueriesTest(t)
	user := createTestUser(t, q)

	fromAccount := createTestAccount(t, q, user.ID, 10000)
	toAccount := createTestAccount(t, q, user.ID, 10000)

	const numGoroutines = 10
	const transferAmount int64 = 100

	var wg sync.WaitGroup
	var errCount atomic.Int64

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			_, err := q.CreateEntry(context.Background(), db.CreateEntryParams{
				AccountID: fromAccount.ID,
				Amount:    -transferAmount,
				Type:      "transfer",
				Category:  "transfer",
			})
			if err != nil {
				errCount.Add(1)
				return
			}

			_, err = q.CreateEntry(context.Background(), db.CreateEntryParams{
				AccountID: toAccount.ID,
				Amount:    transferAmount,
				Type:      "transfer",
				Category:  "transfer",
			})
			if err != nil {
				errCount.Add(1)
				return
			}

			_, err = q.DebitAccount(context.Background(), db.DebitAccountParams{
				Balance: -transferAmount,
				ID:      fromAccount.ID,
				UserID:  user.ID,
			})
			if err != nil {
				errCount.Add(1)
				return
			}

			err = q.CreditAccount(context.Background(), db.CreditAccountParams{
				Balance: transferAmount,
				ID:      toAccount.ID,
				UserID:  user.ID,
			})
			if err != nil {
				errCount.Add(1)
			}
		}()
	}

	wg.Wait()
	require.Zero(t, errCount.Load(), "concurrent transfers produced errors")

	fromResult, err := q.GetAccountByID(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	toResult, err := q.GetAccountByID(context.Background(), toAccount.ID)
	require.NoError(t, err)

	totalExpected := int64(20000)
	totalActual := fromResult.Balance + toResult.Balance
	require.Equal(t, totalExpected, totalActual, "conservation of money violated: total balance changed")

	t.Logf("from balance: %d, to balance: %d, total: %d", fromResult.Balance, toResult.Balance, totalActual)
}

func TestTransferOverdraftRejection(t *testing.T) {
	q := TxQueriesTest(t)
	user := createTestUser(t, q)

	account := createTestAccount(t, q, user.ID, 500)
	targetAccount := createTestAccount(t, q, user.ID, 0)

	const transferAmount int64 = 1000

	rows, err := q.DebitAccount(context.Background(), db.DebitAccountParams{
		Balance: -transferAmount,
		ID:      account.ID,
		UserID:  user.ID,
	})
	require.NoError(t, err)
	require.Zerof(t, rows, "expected 0 rows affected (overdraft rejected), got %d", rows)

	afterDebit, err := q.GetAccountByID(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, int64(500), afterDebit.Balance, "balance should remain unchanged after rejected debit")

	_ = targetAccount
}

func TestTransferMultipleEntriesConsistency(t *testing.T) {
	q := TxQueriesTest(t)
	user := createTestUser(t, q)

	fromAccount := createTestAccount(t, q, user.ID, 50000)
	toAccount := createTestAccount(t, q, user.ID, 0)

	const numTransfers = 20
	const transferAmount int64 = 100

	for i := 0; i < numTransfers; i++ {
		_, err := q.CreateEntry(context.Background(), db.CreateEntryParams{
			AccountID: fromAccount.ID,
			Amount:    -transferAmount,
			Type:      "transfer",
			Category:  fmt.Sprintf("transfer-%d", i),
		})
		require.NoError(t, err)

		_, err = q.CreateEntry(context.Background(), db.CreateEntryParams{
			AccountID: toAccount.ID,
			Amount:    transferAmount,
			Type:      "transfer",
			Category:  fmt.Sprintf("transfer-%d", i),
		})
		require.NoError(t, err)

		_, err = q.DebitAccount(context.Background(), db.DebitAccountParams{
			Balance: -transferAmount,
			ID:      fromAccount.ID,
			UserID:  user.ID,
		})
		require.NoError(t, err)

		err = q.CreditAccount(context.Background(), db.CreditAccountParams{
			Balance: transferAmount,
			ID:      toAccount.ID,
			UserID:  user.ID,
		})
		require.NoError(t, err)
	}

	fromResult, err := q.GetAccountByID(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	toResult, err := q.GetAccountByID(context.Background(), toAccount.ID)
	require.NoError(t, err)

	expectedFrom := int64(50000 - numTransfers*transferAmount)
	expectedTo := int64(numTransfers * transferAmount)

	require.Equal(t, expectedFrom, fromResult.Balance, "source account balance mismatch")
	require.Equal(t, expectedTo, toResult.Balance, "destination account balance mismatch")
	require.Equal(t, int64(50000), fromResult.Balance+toResult.Balance, "conservation of money violated")
}
