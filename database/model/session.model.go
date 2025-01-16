package model

import (
	"time"
)

type Session struct {
	Model

	Users []*User `gorm:"many2many:user_sessions;foreignKey:Id;references:Id;joinForeignKey:SessionId;joinReferences:UserId;constraint:OnDelete:CASCADE" json:"users"`

	RefreshToken string `gorm:"type:varchar(512);not null" json:"refresh_token"`
	IsRevoked    bool   `gorm:"default:false" json:"is_revoked"`

	UserAgent string `gorm:"type:text;not null" json:"user_agent"`
	IPAddress string `gorm:"type:text;not null" json:"ip_address"`

	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}
