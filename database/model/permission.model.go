package model

type Permission struct {
	Model

	Resource string `gorm:"type:varchar(255);not null"`
	Action   string `gorm:"type:varchar(50);not null"`
}
