package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julianinsua/the_simp_bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(33))
	require.NoError(t, err)

	username := util.RandomOwner()

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims)
	require.NotZero(t, claims.Payload.ID)
	require.Equal(t, claims.Username, username)
	require.WithinDuration(t, issuedAt, claims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, claims.ExpiresAt.Time, time.Second)
}

func TestJWTMakerExpired(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(33))
	require.NoError(t, err)

	username := util.RandomOwner()
	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestJWTMakerNoneAlgorithm(t *testing.T) {
	payload, err := NewJWTPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	maker, err := NewJWTMaker(util.RandomString(33))
	require.NoError(t, err)
	payload, err = maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
