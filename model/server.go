package model

import (
	"fmt"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"strings"
	"time"
)

// Server structure
type Server struct {
	Name                string    `json:"name"`
	Address             string    `json:"address"`
	ListenPort          int       `json:"listenPort"`
	PrivateKey          string    `json:"privateKey"`
	PublicKey           string    `json:"publicKey"`
	PresharedKey        string    `json:"presharedKey"`
	Endpoint            string    `json:"endpoint"`
	PersistentKeepalive int       `json:"persistentKeepalive"`
	Dns                 string    `json:"dns"`
	PreUp               string    `json:"preUp"`
	PostUp              string    `json:"postUp"`
	PreDown             string    `json:"preDown"`
	PostDown            string    `json:"postDown"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
}

func (a Server) IsValid() []error {
	errs := make([]error, 0)

	// check if the name empty
	if a.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	// check the name field is between 3 to 40 chars
	if len(a.Name) < 2 || len(a.Name) > 40 {
		errs = append(errs, fmt.Errorf("name must be between 2-40 chars"))
	}
	// check if the address empty
	if a.Address == "" {
		errs = append(errs, fmt.Errorf("address is required"))
	}
	// check if the address are valid
	for _, address := range strings.Split(a.Address, ",") {
		if !util.IsValidCidr(strings.TrimSpace(address)) {
			errs = append(errs, fmt.Errorf("address %s is invalid", address))
		}
	}
	// check if the listenPort is valid
	if a.ListenPort < 0 || a.ListenPort > 65535 {
		errs = append(errs, fmt.Errorf("listenPort %s is invalid", a.ListenPort))
	}
	// check if the endpoint empty
	if a.Endpoint == "" {
		errs = append(errs, fmt.Errorf("endpoint is required"))
	}
	// check if the persistentKeepalive is valid
	if a.PersistentKeepalive < 0 {
		errs = append(errs, fmt.Errorf("persistentKeepalive %d is invalid", a.PersistentKeepalive))
	}
	// check if the dns empty
	if a.Dns == "" {
		errs = append(errs, fmt.Errorf("dns is required"))
	}
	// check if the address are valid
	for _, dns := range strings.Split(a.Dns, ",") {
		if !util.IsValidIp(strings.TrimSpace(dns)) {
			errs = append(errs, fmt.Errorf("dns %s is invalid", dns))
		}
	}

	return errs
}
