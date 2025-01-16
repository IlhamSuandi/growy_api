package model

type Role struct {
	Model

	BranchId uint  `gorm:"index" json:"branch_id"`
	UserId   *uint `gorm:"uniqueIndex" json:"user_id"`

	Name string `gorm:"type:varchar(100);not null" json:"name"`

	// Permissions []*Permission `gorm:"foreignkey:RoleId;references:Id;constraint:OnDelete:SET NULL" json:"permissions"`
}
