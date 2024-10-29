package model

type Log struct {
	Model
	RemoteAddr    *string `gorm:"varchar(50)"`
	Action        *string `gorm:"varchar(10)"`
	Method        *string `gorm:"varchar(10)"`
	Path          *string `gorm:"varchar(50)"`
	Status        uint    `gorm:"type:uint;not null"`
	ExecutionTime string  `gorm:"type:varchar(50)"`
	Size          int     `gorm:"type:int"`
	UserAgent     string  `gorm:"type:text"`
}
