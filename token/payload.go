package token

import (
	"time"

	"github.com/google/uuid"
)

// Holds the payload that is encoded inside the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Creates a new token payload using a username and a duration.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{tokenId, username, time.Now(), time.Now().Add(duration)}
	return &payload, nil
}
