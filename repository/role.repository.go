package repository

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role *gorm.DB) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (rr *roleRepository) CreateRole(role *gorm.DB) error {
	result := rr.db.Create(&role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("can't create role")
	}

	return nil
}

func (rr *roleRepository) AddRolePermission(roleId uint, permission model.Permission) error {
	role := model.Role{
		Model: model.Model{
			Id: roleId,
		},
		Permissions: []*model.Permission{
			&permission,
		},
	}

	result := rr.db.Save(&role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("can't add role permission")
	}

	return nil
}

// func (rr *roleRepository) GetRolePermissions(roleName string) ([]model.Permission, error) {
//   var permissions []model.Permission
//   result := rr.db.Find(&permissions)
// }
