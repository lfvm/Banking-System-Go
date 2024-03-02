package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {

	password := RandomString(6)

	hashedPwd, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd)

	err = CheckPassword(password, hashedPwd)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPwd)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
