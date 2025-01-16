package model

import (
	"time"

	"github.com/google/uuid"
)

type CompanyOption struct {
	Model

	CompanyId    uint      `gorm:"uniqueIndex" json:"company_id"`
	CheckInTime  time.Time `gorm:"type:time;not null" json:"check_in_time"`
	WorkingHours uint      `gorm:"type:integer;not null" jjson:"working_hours"`

	UseCheckout  bool      `gorm:"type:boolean;default:true" json:"use_checkout"`
  CheckOutTime time.Time `gorm:"type:time" json:"check_out_time"`
}

type UserOption struct {
	Id       uint      `gorm:"primaryKey;autoIncrement"`
	OptionId uuid.UUID `gorm:"type:uuid;unique;default:gen_random_uuid()"`

	// TODO: add user options
}
