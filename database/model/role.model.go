package model

type Role struct {
	Model

	BranchId uint
	UserId   *uint `gorm:"uniqueIndex"`

	Name        string        `gorm:"type:varchar(100);not null"`
	Permissions []*Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}
