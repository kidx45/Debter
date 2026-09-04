package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestUserDirect(t *testing.T) db.User {
	var userID int64
	err := testDB.QueryRowContext(context.Background(),
		`INSERT INTO users (username, full_name, email, hashed_password)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		util.RandomString(6), util.RandomString(6), util.RandomString(6), util.RandomString(6),
	).Scan(&userID)
	require.NoError(t, err)
	return db.User{ID: userID}
}

func createTestAccountDirect(t *testing.T, userID int64, accountType string, initialBalance int64) db.Account {
	var accountID int64
	err := testDB.QueryRowContext(context.Background(),
		`INSERT INTO accounts (user_id, account_type, account_number, balance)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		userID, accountType, util.RandomNumber(100000, 999999), initialBalance,
	).Scan(&accountID)
	require.NoError(t, err)
	return db.Account{ID: accountID, UserID: userID, AccountType: accountType, Balance: initialBalance}
}

func TestCreateAccount(t *testing.T) {
	user, q := CreateRandomUser(t)

	account, err := q.CreateAccount(context.Background(), db.CreateAccountParams{
		UserID:        user.ID,
		AccountType:   "savings",
		AccountNumber: util.RandomNumber(100000, 999999),
		Balance:       5000,
	})
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, user.ID, account.UserID)
	require.Equal(t, "savings", account.AccountType)
	require.Equal(t, int64(5000), account.Balance)
}

func TestCreateEntry(t *testing.T) {
	user, q := CreateRandomUser(t)

	account, err := q.CreateAccount(context.Background(), db.CreateAccountParams{
		UserID:        user.ID,
		AccountType:   "checking",
		AccountNumber: util.RandomNumber(100000, 999999),
		Balance:       0,
	})
	require.NoError(t, err)

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

func TestTransactionEntryAndBalanceUpdate(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 10000)

	const entryAmount int64 = 2500

	err := util.ExecTx(context.Background(), testDB, func(q *db.Queries) error {
		_, txErr := q.CreateEntry(context.Background(), db.CreateEntryParams{
			AccountID: account.ID,
			Amount:    -entryAmount,
			Type:      "expense",
			Category:  "food",
		})
		if txErr != nil {
			return txErr
		}

		rows, txErr := q.DebitAccount(context.Background(), db.DebitAccountParams{
			Balance: -entryAmount,
			ID:      account.ID,
			UserID:  user.ID,
		})
		if txErr != nil {
			return txErr
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})
	require.NoError(t, err)

	var balance int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(7500), balance, "balance should be 10000 - 2500 after transaction")

	var entryCount int
	err = testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM entries WHERE account_id = $1`, account.ID,
	).Scan(&entryCount)
	require.NoError(t, err)
	require.Equal(t, 1, entryCount, "one entry should exist after transaction")
}

func TestTransactionRollbackOnInsufficientBalance(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 500)

	const entryAmount int64 = 1000

	err := util.ExecTx(context.Background(), testDB, func(q *db.Queries) error {
		_, txErr := q.CreateEntry(context.Background(), db.CreateEntryParams{
			AccountID: account.ID,
			Amount:    -entryAmount,
			Type:      "expense",
			Category:  "food",
		})
		if txErr != nil {
			return txErr
		}

		rows, txErr := q.DebitAccount(context.Background(), db.DebitAccountParams{
			Balance: -entryAmount,
			ID:      account.ID,
			UserID:  user.ID,
		})
		if txErr != nil {
			return txErr
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})
	require.ErrorIs(t, err, sql.ErrNoRows)

	var balance int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(500), balance, "balance should remain 500 after rollback")

	var entryCount int
	err = testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM entries WHERE account_id = $1`, account.ID,
	).Scan(&entryCount)
	require.NoError(t, err)
	require.Equal(t, 0, entryCount, "no entries should exist after rollback")
}

func TestTransactionConcurrentEntryAndBalanceUpdate(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 50000)

	const numGoroutines = 20
	const entryAmount int64 = 100

	var wg sync.WaitGroup
	var errCount atomic.Int64

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			txErr := util.ExecTx(context.Background(), testDB, func(q *db.Queries) error {
				_, txErr := q.CreateEntry(context.Background(), db.CreateEntryParams{
					AccountID: account.ID,
					Amount:    -entryAmount,
					Type:      "expense",
					Category:  fmt.Sprintf("category-%d", idx),
				})
				if txErr != nil {
					return txErr
				}

				rows, txErr := q.DebitAccount(context.Background(), db.DebitAccountParams{
					Balance: -entryAmount,
					ID:      account.ID,
					UserID:  user.ID,
				})
				if txErr != nil {
					return txErr
				}
				if rows == 0 {
					return sql.ErrNoRows
				}

				return nil
			})
			if txErr != nil {
				errCount.Add(1)
			}
		}(i)
	}

	wg.Wait()

	var balance int64
	err := testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)

	var entryCount int
	err = testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM entries WHERE account_id = $1`, account.ID,
	).Scan(&entryCount)
	require.NoError(t, err)

	require.Equal(t, balance, 50000-int64(entryCount)*entryAmount,
		"balance must equal initial balance minus total debited across all committed entries")

	t.Logf("committed entries: %d, final balance: %d", entryCount, balance)
}
