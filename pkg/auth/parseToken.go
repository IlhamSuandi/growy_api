package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/ilhamSuandi/business_assistant/config"
	types "github.com/ilhamSuandi/business_assistant/types"
)

func ParseToken(tokenString string) (types.JwtClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &types.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWT_SECRET), nil
	})
	// Check if the error is related to the token expiration
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return types.JwtClaims{}, fmt.Errorf("token is expired: %w", err)
			}
		}
		return types.JwtClaims{}, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*types.JwtClaims)
	if !ok {
		return types.JwtClaims{}, fmt.Errorf("invalid token claims %w", err)
	}

	return *claims, nil
}
