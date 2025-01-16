package model

import (
	"time"
)

type QRCode struct {
	Model

	UserId uint   `gorm:"uniqueIndex" json:"user_id"`
	Code   string `gorm:"type:text;not null" json:"code"`
	IsUsed bool   `gorm:"default:false" json:"is_used"`

	ExpiresAt *time.Time `gorm:"type:timestamp" json:"expires_at"`
}
