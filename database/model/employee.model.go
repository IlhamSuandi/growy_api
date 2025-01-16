package model

type Employee struct {
	Model

	Pending       bool          `gorm:"type:boolean;default:false" json:"pending"`
	BranchId      uint          `gorm:"index" json:"branch_id"`
	EmployeeEmail string        `gorm:"uniqueIndex" json:"employee_email"`
	Position      string        `gorm:"type:varchar(255)" json:"position"`
	Salary        *Salary       `gorm:"foreignKey:EmployeeEmail;references:EmployeeEmail;constraint:OnDelete:CASCADE" json:"salary"`
	WorkSchedule  *WorkSchedule `gorm:"foreignKey:EmployeeId;references:Id;constraint:OnDelete:CASCADE" json:"work_schedule"`
}
