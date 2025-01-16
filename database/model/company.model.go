package model

type Company struct {
	Model

	Name       string        `gorm:"types:varchar(255);not null" json:"name"`
	Address    string        `gorm:"types:text;not null" json:"address"`
	OwnerEmail string        `gorm:"index;" json:"owner_email"`
  Options    CompanyOption `gorm:"foreignKey:CompanyId;references:Id;constraint:OnDelete:CASCADE" json:"options"`

	Branches     []Branch      `gorm:"foreignKey:CompanyId;references:Id;constraint:OnDelete:CASCADE" json:"branches"`
	WorkSchedule *WorkSchedule `gorm:"foreignKey:CompanyId;references:Id;constraint:OnDelete:CASCADE" json:"work_schedule"`
}
