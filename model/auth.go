package model

// Auth structure
type Auth struct {
	Oauth2   bool   `json:"oauth2"`
	ClientId string `json:"clientId"`
	Code     string `json:"code"`
	State    string `json:"state"`
	CodeUrl  string `json:"codeUrl"`
}
