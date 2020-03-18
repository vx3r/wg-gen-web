package model

import (
	"fmt"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"time"
)

// Server structure
type Server struct {
	Address             []string  `json:"address"`
	ListenPort          int       `json:"listenPort"`
	Mtu                 int       `json:"mtu"`
	PrivateKey          string    `json:"privateKey"`
	PublicKey           string    `json:"publicKey"`
	Endpoint            string    `json:"endpoint"`
	PersistentKeepalive int       `json:"persistentKeepalive"`
	Dns                 []string  `json:"dns"`
	PreUp               string    `json:"preUp"`
	PostUp              string    `json:"postUp"`
	PreDown             string    `json:"preDown"`
	PostDown            string    `json:"postDown"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
}

func (a Server) IsValid() []error {
	errs := make([]error, 0)

	// check if the address empty
	if len(a.Address) == 0 {
		errs = append(errs, fmt.Errorf("address is required"))
	}
	// check if the address are valid
	for _, address := range a.Address {
		if !util.IsValidCidr(address) {
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
	// check if the mtu is valid
	if a.Mtu < 0 {
		errs = append(errs, fmt.Errorf("MTU %d is invalid", a.PersistentKeepalive))
	}
	// check if the address are valid
	for _, dns := range a.Dns {
		if !util.IsValidIp(dns) {
			errs = append(errs, fmt.Errorf("dns %s is invalid", dns))
		}
	}

	return errs
}
