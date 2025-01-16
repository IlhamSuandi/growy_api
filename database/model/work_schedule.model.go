package model

import "time"

type WorkSchedule struct {
	Model
	WorkingDay string    `gorm:"type:varchar(256)" json:"working_day"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`

	CompanyId  *uint `gorm:"uniqueIndex" json:"company_id"`
	BranchId   *uint `gorm:"uniqueIndex" json:"branch_id"`
	EmployeeId *uint `gorm:"uniqueIndex" json:"employee_id"`
}
