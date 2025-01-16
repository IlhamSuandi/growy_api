package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
)

var (
	AdminRole    model.Role
	EmployeeRole model.Role
)

// func SeedRole(db *gorm.DB) {
// 	log := utils.Log
// 	log.Info("seeding role")
//
// 	AdminRole = model.Role{
// 		UserId:   &UserAdmin.Id,
// 		Name:     "admin",
// 		BranchId: MainBranch.Id,
// 		Permissions: []*model.Permission{
// 			{
// 				Resource: "*",
// 				Action:   "all",
// 			},
// 		},
// 	}
//
// 	EmployeeRole = model.Role{
// 		UserId:   &UserOne.Id,
// 		Name:     "employee",
// 		BranchId: MainBranch.Id,
// 		Permissions: []*model.Permission{
// 			{
// 				Resource: "attendances",
// 				Action:   "get,post",
// 			},
//
// 			{
// 				Resource: "qrcode",
// 				Action:   "get",
// 			},
// 		},
// 	}
//
// 	roles := []*model.Role{
// 		&AdminRole,
// 		&EmployeeRole,
// 	}
//
// 	result := db.CreateInBatches(&roles, 2)
// 	if result.Error != nil {
// 		log.Fatal(result.Error)
// 	}
//
// 	if result.RowsAffected == 0 {
// 		log.Fatal("no rows affected")
// 	}
//
// 	// update user role
// 	if err := db.Save(&UserAdmin).Error; err != nil {
// 		log.Fatal(err)
// 	}
//
// 	if err := db.Save(&UserOne).Error; err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Info("successfully seeded role")
// }
