package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)

// The basic struct that holds the tokens related to the claims in a token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// Holds the payload that is encoded inside the JWT token
type JWTPayload struct {
	Payload
	jwt.RegisteredClaims
}

// Creates a new JWT token payload using a username and a duration.
func NewJWTPayload(username string, duration time.Duration) (*JWTPayload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := JWTPayload{Payload{tokenId, username}, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "the_simp_bank",
		Subject:   username,
		ID:        tokenId.String(),
	}}
	return &payload, nil
}

// Holds the payload that is encoded inside the PASETO token
type PASETOPayload struct {
	Payload
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"IssuedAt"`
}

// Creates a new PASETO token payload using username and duration.
func NewPASETOPayload(username string, duration time.Duration) (*PASETOPayload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &PASETOPayload{Payload{tokenId, username}, time.Now().Add(duration), time.Now()}
	return payload, nil
}
