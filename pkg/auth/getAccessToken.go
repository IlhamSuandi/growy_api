package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAccessToken(r *http.Request) (string, error) {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		accessToken := strings.Split(tokenAuth, "Bearer ")[1]
		return accessToken, nil
	}

	return "", errors.New("token not found")
}
