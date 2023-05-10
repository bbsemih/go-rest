package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPass, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	err = CheckPassword(password, hashedPass)
	require.NoError(t, err)

	wrongPassword := RandomString(7)
	err = CheckPassword(wrongPassword, hashedPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPass2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, hashedPass, hashedPass2)
}
