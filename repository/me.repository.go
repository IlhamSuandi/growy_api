package repository

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type MeRepository interface {
	GetMe(userId uint) (*model.User, error)
}

type meRepository struct {
	db *gorm.DB
}

func NewMeRepository(db *gorm.DB) MeRepository {
	return &meRepository{
		db: db,
	}
}

func (mr *meRepository) GetMe(userId uint) (*model.User, error) {
	var user model.User
	result := mr.db.Where("id = ?", userId).Preload("Employee").Preload("Employee.Salary").Preload("Employee.WorkSchedule").Preload("Attendances").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
