package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	// Test right password
	pswd := RandomString(6)
	hash1, err := HashPassword(pswd)

	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	err = CheckPassword(pswd, hash1)

	require.NoError(t, err)

	// Test wrong password
	wrong := RandomString(6)
	err = CheckPassword(wrong, hash1)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hash2, err := HashPassword(pswd)

	require.NoError(t, err)
	require.NotEmpty(t, hash2)

	require.NotEqual(t, hash1, hash2)
}
