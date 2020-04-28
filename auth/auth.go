package auth

import (
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"golang.org/x/oauth2"
)

type Auth interface {
	Setup() error
	CodeUrl(state string) string
	Exchange(code string) (*oauth2.Token, error)
	UserInfo(oauth2Token *oauth2.Token) (*model.User, error)
}
