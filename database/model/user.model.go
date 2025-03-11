package model

type User struct {
	Model
	Username        string `gorm:"type:varchar(100)" json:"username"`
	Email           string `gorm:"type:varchar(50);unique;not null" json:"email"`
	Password        string `gorm:"type:varchar(256);not null" json:"-"`
	IsEmailVerified bool   `gorm:"type:boolean;default:false" json:"is_email_verified"`
	AuthProvider    string `gorm:"type:varchar(50)" json:"auth_provider"`
	IsOnBoarded     bool   `gorm:"type:boolean;default:false" json:"is_on_boarded"`
	Role            string `gorm:"type:varchar(10);not null" json:"role"` // admin, owner, employee
	Picture         string `gorm:"type:varchar(255)" json:"picture"`

	Log         *Log          `gorm:"foreignKey:Email;references:Email;constraint:OnDelete:CASCADE"`
	Company     []Company     `gorm:"foreignKey:OwnerEmail;references:Email;constraint:OnDelete:CASCADE" json:"company"`
	Employee    *Employee     `gorm:"foreignKey:EmployeeEmail;references:Email;constraint:OnDelete:CASCADE" json:"employee"`
	Permissions []*Permission `gorm:"foreignkey:UserId;references:Id;constraint:OnDelete:SET NULL" json:"permissions"`

	QRCode *QRCode `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`

	Sessions    []*Session    `gorm:"many2many:user_sessions;foreignKey:Id;references:Id;joinForeignKey:UserId;joinReferences:SessionId;constraint:OnDelete:CASCADE"`
	Attendances []*Attendance `gorm:"many2many:user_attendances;foreignKey:Id;references:Id;joinForeignKey:UserId;joinReferences:AttendanceId;constraint:OnDelete:CASCADE"`
}
