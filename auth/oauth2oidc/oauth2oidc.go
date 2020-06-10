package oauth2oidc

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"golang.org/x/oauth2"
	"os"
)

// Oauth2idc in order to implement interface, struct is required
type Oauth2idc struct{}

var (
	oauth2Config        *oauth2.Config
	oidcProvider        *oidc.Provider
	oidcIDTokenVerifier *oidc.IDTokenVerifier
)

// Setup validate provider
func (o *Oauth2idc) Setup() error {
	var err error

	oidcProvider, err = oidc.NewProvider(context.TODO(), os.Getenv("OAUTH2_PROVIDER"))
	if err != nil {
		return err
	}

	oidcIDTokenVerifier = oidcProvider.Verifier(&oidc.Config{
		ClientID: os.Getenv("OAUTH2_CLIENT_ID"),
	})

	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH2_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     oidcProvider.Endpoint(),
	}

	return nil
}

// CodeUrl get url to redirect client for auth
func (o *Oauth2idc) CodeUrl(state string) string {
	return oauth2Config.AuthCodeURL(state)
}

// Exchange exchange code for Oauth2 token
func (o *Oauth2idc) Exchange(code string) (*oauth2.Token, error) {
	oauth2Token, err := oauth2Config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	return oauth2Token, nil
}

// UserInfo get token user
func (o *Oauth2idc) UserInfo(oauth2Token *oauth2.Token) (*model.User, error) {
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	iDToken, err := oidcIDTokenVerifier.Verify(context.TODO(), rawIDToken)
	if err != nil {
		return nil, err
	}

	userInfo, err := oidcProvider.UserInfo(context.TODO(), oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return nil, err
	}

	// ID Token payload is just JSON
	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to get id token claims: %s", err)
	}

	// get some infos about user
	user := &model.User{}
	user.Sub = userInfo.Subject
	user.Email = userInfo.Email
	user.Profile = userInfo.Profile

	if v, found := claims["name"]; found && v != nil {
		user.Name = v.(string)
	} else {
		log.Error("name not found in user info claims")
	}

	user.Issuer = iDToken.Issuer
	user.IssuedAt = iDToken.IssuedAt

	return user, nil
}
