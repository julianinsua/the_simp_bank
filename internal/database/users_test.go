package database

import (
	"context"
	"testing"
	"time"

	"github.com/julianinsua/the_simp_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hash, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	params := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hash,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	usr, err := testQueries.CreateUser(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, usr)
	require.Equal(t, usr.Username, params.Username)
	require.Equal(t, usr.HashedPassword, params.HashedPassword)
	require.Equal(t, usr.FullName, params.FullName)
	require.Equal(t, usr.Email, params.Email)

	require.True(t, usr.PasswordChangedAt.IsZero())
	require.NotZero(t, usr.CreatedAt)

	return usr
}

func TestUserCreation(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
