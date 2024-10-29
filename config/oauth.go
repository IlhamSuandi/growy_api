package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func init() {
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     OAUTH2_CLIENT_ID,
		ClientSecret: OAUTH2_CLIENT_SECRET,
		RedirectURL:  OAUTH2_ClIENT_REDIRECT_URL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
