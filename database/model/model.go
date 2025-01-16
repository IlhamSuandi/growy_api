package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	UUID      uuid.UUID `gorm:"type:uuid;unique;default:gen_random_uuid()" json:"uuid"`
	CreatedAt time.Time `gorm:"autoCreateTime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
