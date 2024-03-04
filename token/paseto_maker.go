package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {

	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symectric key must be of size %d, have %d", chacha20poly1305.KeySize, len(symetricKey))
	}

	return &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}, nil

}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	paylaod, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symetricKey, paylaod, nil)

}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symetricKey, payload, nil)

	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()

	if err != nil {
		return nil, ErrExpiredToken
	}

	return payload, nil

}
