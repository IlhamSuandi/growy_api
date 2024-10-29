package model

type User struct {
	Model
	Username        string `gorm:"type:varchar(100)"`
	Email           string `gorm:"type:varchar(50);unique;not null"`
	Password        string `gorm:"type:varchar(256);not null" json:"-"`
	IsEmailVerified bool   `gorm:"type:boolean;default:false"`

	Company []Company `gorm:"foreignKey:OwnerEmail;references:Email;constraint:OnDelete:CASCADE"`

	AuthProvider string `gorm:"type:varchar(50)"`

	QRCode QRCode `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	Role   Role   `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`

	Sessions    []*Session    `gorm:"many2many:user_sessions;foreignKey:Id;references:Id;joinForeignKey:UserId;joinReferences:SessionId;constraint:OnDelete:CASCADE"`
	Attendances []*Attendance `gorm:"many2many:user_attendances;foreignKey:Id;references:Id;joinForeignKey:UserId;joinReferences:AttendanceId;constraint:OnDelete:CASCADE"`
}
