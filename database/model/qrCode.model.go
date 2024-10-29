package model

import (
	"time"
)

type QRCode struct {
	Model

	UserId uint   `gorm:"uniqueIndex"`
	Code   string `gorm:"type:text;not null"`
	IsUsed bool   `gorm:"default:false"`

	ExpiresAt *time.Time `gorm:"type:timestamp"`
}
