package model

type Branch struct {
	Model
	CompanyId uint `gorm:"not null;index:idx_company_name,unique" json:"company_id"`

	Name    string `gorm:"type:varchar(255);not null;index:idx_company_name,unique" json:"name"`
	Address string `gorm:"types:text;not null" json:"address"`

	Roles        []Role        `gorm:"foreignKey:BranchId;references:Id;constraint:OnDelete:CASCADE" json:"roles"`
	Employees    []Employee    `gorm:"foreignKey:BranchId;references:Id;constraint:OnDelete:CASCADE" json:"employees"`
	WorkSchedule *WorkSchedule `gorm:"foreignKey:BranchId;references:Id;constraint:OnDelete:CASCADE" json:"work_schedule"`
}
