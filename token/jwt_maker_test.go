package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lfvm/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {

	maker, err := NewJwtMaker(utils.RandomString(32))
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

func TestExpiredJwt(t *testing.T) {
	maker, err := NewJwtMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()

	jwt, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)

	payload, err := maker.VerifyToken(jwt)
	require.Error(t, err)
	require.Empty(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewJwtMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	paylaod, err := NewPayload(username, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, paylaod)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
}
