package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
)

const minSymmetricKeySize = 32

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != minSymmetricKeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", minSymmetricKeySize)
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(userID int64, username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
