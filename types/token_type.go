package types

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserId    uuid.UUID
	SessionId uuid.UUID
	Username  string
	Email     string
	Role      string
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken            string
	AccessTokenExpiration  time.Duration
	RefreshToken           string
	RefreshTokenExpiration time.Duration
	SessionId              uuid.UUID
}

type CreateToken struct {
	UserId    uuid.UUID
	SessionId uuid.UUID
	Username  string
	Email     string
	Duration  time.Duration
}
