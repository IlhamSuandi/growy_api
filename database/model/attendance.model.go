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
	CheckIn  *time.Time `gorm:"default:null" json:"check_in"`
	CheckOut *time.Time `gorm:"default:null" json:"check_out"`
	Status   string     `gorm:"type:varchar(20);not null" json:"status"`
	Date     time.Time  `gorm:"type:date;not null" json:"date"`
	Location *string    `gorm:"type:varchar(100);default:null" json:"location"`

  Users []*User `gorm:"many2many:user_attendances;foreignKey:Id;references:Id;joinForeignKey:AttendanceId;joinReferences:UserId;constraint:OnDelete:CASCADE" json:"users"`
}
