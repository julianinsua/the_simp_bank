package token

import "time"

// Interface to manage tokens
type Maker interface {
	// Creates a new tokenfor a specific user and for a given duration
	CreateToken(username string, duration time.Duration) (string, error)
	// Check if a given token is valid or not
	VerifyToken(token string) (*Payload, error)
}
