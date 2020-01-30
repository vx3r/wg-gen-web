package model

import "time"

type Client struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Enable  bool      `json:"enable"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	AllowedIPs string `json:"allowedIPs"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
