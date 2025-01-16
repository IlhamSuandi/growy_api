package model

type Permission struct {
	Model
	UserId uint `gorm:"not null;uniqueIndex:idx_user_resource" json:"user_id"`

	Resource string `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_resource" json:"resource"`
	Action   string `gorm:"type:varchar(50);not null" json:"action"` // get,post,put/patch,delete
}
