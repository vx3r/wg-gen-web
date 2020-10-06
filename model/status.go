package model

import (
	"time"
)

// ClientStatus structure
type ClientStatus struct {
	PublicKey        string    `json:"publicKey"`
	HasPresharedKey  bool      `json:"hasPresharedKey"`
	ProtocolVersion  int       `json:"protocolVersion"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Connected        bool      `json:"connected"`
	AllowedIPs       []string  `json:"allowedIPs"`
	Endpoint         string    `json:"endpoint"`
	LastHandshake    time.Time `json:"lastHandshake"`
	ReceivedBytes    int       `json:"receivedBytes"`
	TransmittedBytes int       `json:"transmittedBytes"`
}

// InterfaceStatus structure
type InterfaceStatus struct {
	Name          string `json:"name"`
	DeviceType    string `json:"type"`
	ListenPort    int    `json:"listenPort"`
	NumberOfPeers int    `json:"numPeers"`
	PublicKey     string `json:"publicKey"`
}
