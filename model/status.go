package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// ClientStatus structure
type ClientStatus struct {
	PublicKey             string        `json:"publicKey"`
	HasPresharedKey       bool          `json:"hasPresharedKey"`
	ProtocolVersion       int           `json:"protocolVersion"`
	Name                  string        `json:"name"`
	Email                 string        `json:"email"`
	Connected             bool          `json:"connected"`
	AllowedIPs            []string      `json:"allowedIPs"`
	Endpoint              string        `json:"endpoint"`
	LastHandshake         time.Time     `json:"lastHandshake"`
	LastHandshakeRelative time.Duration `json:"lastHandshakeRelative"`
	ReceivedBytes         int           `json:"receivedBytes"`
	TransmittedBytes      int           `json:"transmittedBytes"`
}

func (c *ClientStatus) MarshalJSON() ([]byte, error) {

	duration := fmt.Sprintf("%v ago", c.LastHandshakeRelative)
	if c.LastHandshakeRelative.Hours() > 5208 { // 24*7*31 = approx one month
		duration = "more than a month ago"
	}
	return json.Marshal(&struct {
		PublicKey             string    `json:"publicKey"`
		HasPresharedKey       bool      `json:"hasPresharedKey"`
		ProtocolVersion       int       `json:"protocolVersion"`
		Name                  string    `json:"name"`
		Email                 string    `json:"email"`
		Connected             bool      `json:"connected"`
		AllowedIPs            []string  `json:"allowedIPs"`
		Endpoint              string    `json:"endpoint"`
		LastHandshake         time.Time `json:"lastHandshake"`
		LastHandshakeRelative string    `json:"lastHandshakeRelative"`
		ReceivedBytes         int       `json:"receivedBytes"`
		TransmittedBytes      int       `json:"transmittedBytes"`
	}{
		PublicKey:             c.PublicKey,
		HasPresharedKey:       c.HasPresharedKey,
		ProtocolVersion:       c.ProtocolVersion,
		Name:                  c.Name,
		Email:                 c.Email,
		Connected:             c.Connected,
		AllowedIPs:            c.AllowedIPs,
		Endpoint:              c.Endpoint,
		LastHandshake:         c.LastHandshake,
		LastHandshakeRelative: duration,
		ReceivedBytes:         c.ReceivedBytes,
		TransmittedBytes:      c.TransmittedBytes,
	})
}

// InterfaceStatus structure
type InterfaceStatus struct {
	Name          string `json:"name"`
	DeviceType    string `json:"type"`
	ListenPort    int    `json:"listenPort"`
	NumberOfPeers int    `json:"numPeers"`
	PublicKey     string `json:"publicKey"`
}
