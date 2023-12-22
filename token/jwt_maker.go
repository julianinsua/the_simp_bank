package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// A JWT token generator
type JWTMaker struct {
	secretKey string
}

// Creates a new JWTMaker
func NewJWTMaker(secretKey string) (JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return JWTMaker{}, fmt.Errorf("invalid key size: secret key must be at least %v characters", minSecretKeySize)
	}
	return JWTMaker{secretKey}, nil
}

// JSON web token maker. Implements the Maker interface.
func (mkr *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewJWTPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("unable to cretae token payload")
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(mkr.secretKey))
}

func (j *JWTMaker) VerifyToken(token string) (*JWTPayload, error) {
	keyFunc := func(tk *jwt.Token) (interface{}, error) {
		// check signing algorithm
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := jwtToken.Claims.(*JWTPayload); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}
