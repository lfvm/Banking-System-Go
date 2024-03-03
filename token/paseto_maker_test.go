package token

import (
	"testing"
	"time"

	"github.com/lfvm/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestCreatePasetoToken(t *testing.T) {

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	jwt, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	payload, err := maker.VerifyToken(jwt)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)

	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()

	jwt, err := maker.CreateToken(username, -time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	payload, err := maker.VerifyToken(jwt)

	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}
