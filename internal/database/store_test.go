package database

import (
	"context"
	"testing"

	"github.com/julianinsua/the_simp_bank.git/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// Initialize Store
	store := NewStore(testDB)

	// Create accounts
	acc1, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)

	acc2, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)

	// Concurrent transactions
	errs := make(chan error)
	txResults := make(chan TransferTxResult)

	amount := 10.0
	n := 5
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			txResults <- result
		}()
	}

	existed := make(map[int]bool)
	// Listen for results and errors
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		res := <-txResults
		require.NotEmpty(t, res)

		// check for transfer
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check for acc1 entries
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check for acc2 entries
		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check from account
		fromAcc := res.FromAccount
		require.NotEmpty(t, fromAcc)
		require.Equal(t, fromAcc.ID, acc1.ID)

		// check to account
		toAcc := res.ToAccount
		require.NotEmpty(t, toAcc)
		require.Equal(t, toAcc.ID, acc2.ID)

		// check balances
		diff1 := acc1.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		// Check consecutive transactions are done only once
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check final balances
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.Equal(t, acc1.Balance-float64(n)*amount, updatedAcc1.Balance)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.Equal(t, acc2.Balance+float64(n)*amount, updatedAcc2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	// Initialize Store
	store := NewStore(testDB)

	// Create accounts
	acc1, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)

	acc2, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)

	// Concurrent transactions
	errs := make(chan error)

	amount := 10.0
	n := 10
	for i := 0; i < n; i++ {
		fromAccId := acc1.ID
		toAccId := acc2.ID

		if i%2 == 0 {
			fromAccId = acc2.ID
			toAccId = acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccId,
				ToAccountID:   toAccId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// Listen for results and errors
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check final balances
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.Equal(t, acc1.Balance, updatedAcc1.Balance)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
