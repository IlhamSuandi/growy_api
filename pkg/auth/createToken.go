package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/types"
)

func CreateToken(session types.CreateToken) (string, *types.JwtClaims, error) {
	// create claims
	claims := types.JwtClaims{
		SessionId: session.SessionId,
		UserId:    session.UserId,
		Username:  session.Username,
		Email:     session.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			Subject:   session.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(session.Duration).Unix(),
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// create token string
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &claims, nil
}
