package auth

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/auth/fake"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/auth/github"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/auth/oauth2oidc"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"golang.org/x/oauth2"
	"os"
)

// Auth interface to implement as auth provider
type Auth interface {
	Setup() error
	CodeUrl(state string) string
	Exchange(code string) (*oauth2.Token, error)
	UserInfo(oauth2Token *oauth2.Token) (*model.User, error)
}

// GetAuthProvider  get an instance of auth provider based on config
func GetAuthProvider() (Auth, error) {
	var oauth2Client Auth
	var err error

	switch os.Getenv("OAUTH2_PROVIDER_NAME") {
	case "fake":
		log.Warn("Oauth is set to fake, no actual authentication will be performed")
		oauth2Client = &fake.Fake{}

	case "oauth2oidc":
		log.Warn("Oauth is set to oauth2oidc, must be RFC implementation on server side")
		oauth2Client = &oauth2oidc.Oauth2idc{}

	case "github":
		log.Warn("Oauth is set to github, no openid will be used")
		oauth2Client = &github.Github{}

	case "google":
		return nil, fmt.Errorf("auth provider name %s not yet implemented", os.Getenv("OAUTH2_PROVIDER_NAME"))
	default:
		return nil, fmt.Errorf("auth provider name %s unknown", os.Getenv("OAUTH2_PROVIDER_NAME"))
	}

	err = oauth2Client.Setup()

	return oauth2Client, err
}
