package model

type Branch struct {
	Model
	CompanyId uint

	Name    string `gorm:"types:varchar(255);not null"`
	Address string `gorm:"types:text;not null"`

	Options BranchOption `gorm:"foreignKey:BranchId;references:Id;constraint:OnDelete:CASCADE"`

	Roles     []Role `gorm:"foreignKey:BranchId;references:Id;constraint:OnDelete:CASCADE"`
	Employees []User `gorm:"many2many:employees;foreignKey:Id;references:Id;joinForeignKey:BranchId;joinReferences:UserId;constraint:OnDelete:CASCADE"`
}
