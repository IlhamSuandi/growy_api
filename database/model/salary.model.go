package model

type Salary struct {
	Model

	EmployeeEmail   string `gorm:"uniqueIndex"`
	Monthly         int    `gorm:"type:integer;default:0"`
	Hourly          int    `gorm:"type:integer;default:0"`
	CurrentEarnings int    `gorm:"type:integer;default:0"`
}
