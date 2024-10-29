package model

import (
	"time"
)

type Session struct {
	Model

	Users []*User `gorm:"many2many:user_sessions;foreignKey:Id;references:Id;joinForeignKey:SessionId;joinReferences:UserId;constraint:OnDelete:CASCADE"`

	RefreshToken string `gorm:"type:varchar(512);not null"`
	IsRevoked    bool   `gorm:"default:false"`

	UserAgent string `gorm:"type:text;not null"`
	IPAddress string `gorm:"type:text;not null"`

	ExpiresAt time.Time `gorm:"not null"`
}
