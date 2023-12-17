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

// Creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: secret key must be at least %v characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// JSON web token maker. Implements the Maker interface.
func (mkr *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("unable to cretae token payload")
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"payload": payload})
	return jwtToken.SignedString(mkr.secretKey)
}

func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {

}
