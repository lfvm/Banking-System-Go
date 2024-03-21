package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {

	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil

}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))

	if err != nil {
		return "", nil, err
	}

	return token, payload, nil

}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	payload := Payload{}

	jwtToken, err := jwt.ParseWithClaims(token, &payload, keyFunc)

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, err
	}

	return &payload, nil

}
