package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	adapter "github.com/kidx45/Debter/internal/adapter/outbound"
	port "github.com/kidx45/Debter/internal/port/outbound"
	"github.com/kidx45/Debter/internal/util"
	"github.com/stretchr/testify/require"
)

func newTransactionRepo() port.UpdateContents {
	return adapter.NewPostgresTransactionRepository(db.New(testDB), testDB)
}

func TestUpdateBalanceTxExpense(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 10000)
	repo := newTransactionRepo()

	const intendedAmount int64 = -2500

	entry, updatedAccount, err := repo.UpdateBalanceTx(context.Background(),
		account.ID, user.ID, intendedAmount, "expense", "food")
	require.NoError(t, err)

	require.Equal(t, intendedAmount, entry.Amount)
	require.Equal(t, "expense", entry.Type)
	require.Equal(t, "food", entry.Category)
	require.Equal(t, account.ID, entry.AccountID)

	require.Equal(t, int64(7500), updatedAccount.Balance,
		"final balance must reflect the intended debit")

	var balance int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(7500), balance, "db balance must match intended end value")

	var entryCount, entryAmountDB int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM entries WHERE account_id = $1`, account.ID,
	).Scan(&entryCount, &entryAmountDB)
	require.NoError(t, err)
	require.Equal(t, int64(1), entryCount)
	require.Equal(t, intendedAmount, entryAmountDB)
}

func TestUpdateBalanceTxIncome(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 1000)
	repo := newTransactionRepo()

	const intendedAmount int64 = 4000

	entry, updatedAccount, err := repo.UpdateBalanceTx(context.Background(),
		account.ID, user.ID, intendedAmount, "income", "salary")
	require.NoError(t, err)

	require.Equal(t, intendedAmount, entry.Amount)
	require.Equal(t, "income", entry.Type)
	require.Equal(t, int64(5000), updatedAccount.Balance,
		"final balance must reflect the intended credit")

	var balance int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(5000), balance)
}

func TestUpdateBalanceTxRollbackOnInsufficientBalance(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 500)
	repo := newTransactionRepo()

	_, _, err := repo.UpdateBalanceTx(context.Background(),
		account.ID, user.ID, -1000, "expense", "food")
	require.ErrorIs(t, err, sql.ErrNoRows)

	var balance int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(500), balance, "balance must be unchanged after rollback")

	var entryCount int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM entries WHERE account_id = $1`, account.ID,
	).Scan(&entryCount)
	require.NoError(t, err)
	require.Equal(t, int64(0), entryCount, "no entry may survive a rolled-back transaction")
}

func TestUpdateBalanceTxConcurrentRandomized(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 100000)
	repo := newTransactionRepo()

	const numGoroutines = 8

	var mu sync.Mutex
	debited := make([]int64, numGoroutines)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			amount := -util.RandomNumber(50, 500)
			_, _, err := repo.UpdateBalanceTx(context.Background(),
				account.ID, user.ID, amount, "expense", fmt.Sprintf("category-%d", idx))
			if err != nil {
				return
			}

			mu.Lock()
			debited[idx] = amount
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	var totalDebited int64
	for _, d := range debited {
		totalDebited += d
	}

	expectedBalance := 100000 + totalDebited

	var balance int64
	err := testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance,
		"balance must equal initial minus total debited across all committed transactions")

	var committed int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM entries WHERE account_id = $1`, account.ID,
	).Scan(&committed)
	require.NoError(t, err)
	require.Equal(t, int64(numGoroutines), committed,
		"each committed transaction must produce exactly one entry")

	t.Logf("committed entries: %d, total debited: %d, final balance: %d", committed, totalDebited, balance)
}

func TestUpdateAccount(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 1000)
	repo := newTransactionRepo()

	err := repo.UpdateAccount(context.Background(), account.ID, user.ID, 2000)
	require.NoError(t, err)

	var balance int64
	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(3000), balance)

	err = repo.UpdateAccount(context.Background(), account.ID, user.ID, -2500)
	require.NoError(t, err)

	err = testDB.QueryRowContext(context.Background(),
		`SELECT balance FROM accounts WHERE id = $1`, account.ID,
	).Scan(&balance)
	require.NoError(t, err)
	require.Equal(t, int64(500), balance)

	err = repo.UpdateAccount(context.Background(), account.ID, user.ID, -1000)
	require.ErrorIs(t, err, sql.ErrNoRows, "overdraw must be rejected")
}

func TestInsertEntry(t *testing.T) {
	user := createTestUserDirect(t)
	account := createTestAccountDirect(t, user.ID, "checking", 0)
	repo := newTransactionRepo()

	entry, err := repo.InsertEntry(context.Background(),
		account.ID, user.ID, 500, "income", "gift")
	require.NoError(t, err)
	require.Equal(t, int64(500), entry.Amount)
	require.Equal(t, "income", entry.Type)
	require.Equal(t, "gift", entry.Category)
	require.Equal(t, account.ID, entry.AccountID)
}