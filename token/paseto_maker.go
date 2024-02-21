package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// A PASETO token generator
type PASETOMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPASETOMaker(symetricKey string) (PASETOMaker, error) {
	if len(symetricKey) != chacha20poly1305.KeySize {
		return PASETOMaker{}, fmt.Errorf("invalid key size: secret key must be at least %d characters", chacha20poly1305.KeySize)
	}
	return PASETOMaker{paseto.NewV2(), []byte(symetricKey)}, nil
}

// Creates a new PASETO V2 symetric token. Implements the Maker Interface
func (mkr *PASETOMaker) CreateToken(username string, duration time.Duration) (string, *PASETOPayload, error) {
	payload, err := NewPASETOPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := mkr.paseto.Encrypt(mkr.symetricKey, payload, nil)
	return token, payload, err
}

// Verifies a token string using PASETO V2 Symetric encoding. Implements the Maker Interface.
func (mkr *PASETOMaker) VerifyToken(token string) (*PASETOPayload, error) {
	payload := &PASETOPayload{}

	err := mkr.paseto.Decrypt(token, mkr.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if payload.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
