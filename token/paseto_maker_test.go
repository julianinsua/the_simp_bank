package token

import (
	"testing"
	"time"

	"github.com/julianinsua/the_simp_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPASETOMaker(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(33))
	require.NoError(t, err)

	username := util.RandomOwner()

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	claims, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims)
	require.NotZero(t, claims.ID)
	require.Equal(t, claims.Username, username)
	require.WithinDuration(t, issuedAt, claims.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, claims.ExpiresAt, time.Second)
}

func TestPASETOMakerExpired(t *testing.T) {
	maker, err := NewPASETOMaker(util.RandomString(33))
	require.NoError(t, err)

	username := util.RandomOwner()
	token, payload, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
