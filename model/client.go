package model

import (
	"fmt"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"strings"
	"time"
)

// Client structure
type Client struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Enable     bool      `json:"enable"`
	AllowedIPs string    `json:"allowedIPs"`
	Address    string    `json:"address"`
	PrivateKey string    `json:"privateKey"`
	PublicKey  string    `json:"publicKey"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

func (a Client) IsValid() []error {
	errs := make([]error, 0)

	// check if the name empty
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	// check the name field is between 3 to 40 chars
	if len(a.Name) < 2 || len(a.Name) > 40 {
		errs = append(errs, fmt.Errorf("name field must be between 2-40 chars"))
	}
	// check if the email empty
	if a.Email == "" {
		errs = append(errs, fmt.Errorf("email field is required"))
	}
	// check if email valid
	if !util.RegexpEmail.MatchString(a.Email) {
		errs = append(errs, fmt.Errorf("email %s is invalid", a.Email))
	}
	// check if the allowedIPs empty
	if a.AllowedIPs == "" {
		errs = append(errs, fmt.Errorf("allowedIPs field is required"))
	}
	// check if the allowedIPs are valid
	for _, allowedIP := range strings.Split(a.AllowedIPs, ",") {
		if !util.IsValidCidr(strings.TrimSpace(allowedIP)) {
			errs = append(errs, fmt.Errorf("allowedIP %s is invalid", allowedIP))
		}
	}
	// check if the address empty
	if a.Address == "" {
		errs = append(errs, fmt.Errorf("address field is required"))
	}
	// check if the address are valid
	for _, address := range strings.Split(a.Address, ",") {
		if !util.IsValidCidr(strings.TrimSpace(address)) {
			errs = append(errs, fmt.Errorf("address %s is invalid", address))
		}
	}

	return errs
}
