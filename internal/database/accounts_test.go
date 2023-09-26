package database

import (
	"context"
	"testing"

	"github.com/julianinsua/the_simp_bank.git/util"
	"github.com/stretchr/testify/require"
)

func TestAccountCreation(t *testing.T) {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, acc.Owner, params.Owner)
	require.Equal(t, acc.Balance, params.Balance)
	require.Equal(t, acc.Currency, params.Currency)

	require.NotEmpty(t, acc.ID)
	require.NotEmpty(t, acc.CreatedAt)
}
