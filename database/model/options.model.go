package model

import "github.com/google/uuid"

type BranchOption struct {
	Model

	BranchId uint `gorm:"uniqueIndex"`

	UseCheckout bool `gorm:"type:boolean;default:true"`
}

type UserOption struct {
	Id       uint      `gorm:"primaryKey;autoIncrement"`
	OptionId uuid.UUID `gorm:"type:uuid;unique;default:gen_random_uuid()"`

	// TODO: add user options
}
