package model

type Company struct {
	Model

	Name       string `gorm:"types:varchar(255);not null"`
	Address    string `gorm:"types:text;not null"`
	OwnerEmail string `gorm:"uniqueIndex;"`

	Branches []Branch `gorm:"foreignKey:CompanyId;references:Id;constraint:OnDelete:CASCADE"`
}
