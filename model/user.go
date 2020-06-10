package model

import "time"

// User structure
type User struct {
	Sub      string    `json:"sub"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Profile  string    `json:"profile"`
	Issuer   string    `json:"issuer"`
	IssuedAt time.Time `json:"issuedAt"`
}
