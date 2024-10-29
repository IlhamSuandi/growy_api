package model

import (
	"time"
)

type AttendanceStatus string

var (
	AttendanceStatusPresent       string = "present"
	AttendanceStatusAbsent        string = "absent"
	AttendanceStatusLate          string = "late"
	AttendanceStatusHalf          string = "half"
	AttendanceStatusApprovedLeave string = "approved_leave"
	AttendanceStatusSick          string = "sick"
)

type Attendance struct {
	Model
	CheckIn  *time.Time `gorm:"default:null"`
	CheckOut *time.Time `gorm:"default:null"`
	Status   string     `gorm:"type:varchar(20);not null"`
	Date     time.Time  `gorm:"type:date;not null"`
	Location *string    `gorm:"type:varchar(100);default:null"`

	Users []*User `gorm:"many2many:user_attendances;foreignKey:Id;references:Id;joinForeignKey:AttendanceId;joinReferences:UserId;constraint:OnDelete:CASCADE"`
}
