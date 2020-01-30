package model

import "time"

// server structure
type Server struct {
	Name                string    `json:"name"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
	Address             string    `json:"address"`
	ListenPort          int       `json:"listenPort"`
	PrivateKey          string    `json:"privateKey"`
	PublicKey           string    `json:"publicKey"`
	PresharedKey        string    `json:"presharedKey"`
	Endpoint            string    `json:"endpoint"`
	PersistentKeepalive int       `json:"persistentKeepalive"`
	Dns                 string    `json:"dns"`
}
