package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID `gorm:"type:uuid;unique;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:CURRENT_TIMESTAMP"`
}
