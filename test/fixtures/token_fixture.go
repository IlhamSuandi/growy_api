package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/pkg/auth"
	"github.com/ilhamSuandi/business_assistant/types"
)

var (
	accessTokenExpiration  = time.Minute * time.Duration(config.JWT_EXPIRATION)
	refreshTokenExpiration = time.Minute * time.Duration(config.JWT_REFRESH_EXPIRATION)
)

func AccessToken(userId uuid.UUID, username string, email string) (string, error) {
	session := types.CreateToken{
		UserId:    userId,
		SessionId: uuid.New(),
		Username:  username,
		Email:     email,
		Duration:  accessTokenExpiration,
	}

	accessToken, _, err := auth.CreateToken(session)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
