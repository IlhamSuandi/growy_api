package seeds

import (
	"errors"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

var (
	AdminRole    model.Role
	EmployeeRole model.Role
)

func SeedRole(db *gorm.DB) error {
	log := utils.Log
	log.Info("seeding role")

	AdminRole = model.Role{
		Name:     "admin",
		BranchId: MainBranch.Id,
		Permissions: []*model.Permission{
			{
				Resource: "*",
				Action:   "all",
			},
		},
	}

	EmployeeRole = model.Role{
		Name:     "employee",
		BranchId: MainBranch.Id,
		Permissions: []*model.Permission{
			{
				Resource: "attendances",
				Action:   "read,create",
			},

			{
				Resource: "qrcode",
				Action:   "read",
			},
		},
	}

	result := db.Create(&AdminRole)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	// update user role
	UserAdmin.Role.Id = AdminRole.Id
	if err := db.Save(&UserAdmin).Error; err != nil {
		return err
	}

	UserOne.Role.Id = EmployeeRole.Id
	if err := db.Save(&UserOne).Error; err != nil {
		return err
	}
	log.Info("successfully seeded role")
	return nil
}
