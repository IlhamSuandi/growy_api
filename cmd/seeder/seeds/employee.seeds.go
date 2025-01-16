package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

func SeedEmployees(db *gorm.DB) {
	log := utils.Log
	log.Info("seeding employees")

	employees := []model.Employee{
		{
			EmployeeEmail: UserOne.Email,
			Position:      "manager",
			BranchId:      MainBranch.Id,
			Salary: &model.Salary{
				Monthly: 10_000_000,
			},
		},
		{
			EmployeeEmail: UserTwo.Email,
			Position:      "chef",
			BranchId:      MainBranch.Id,
			Salary: &model.Salary{
				Monthly: 8_000_000,
			},
		},
		{
			EmployeeEmail: UserThree.Email,
			Position:      "waiter",
			BranchId:      MainBranch.Id,
			Salary: &model.Salary{
				Monthly: 2_500_000,
			},
		},
		{
			EmployeeEmail: UserFour.Email,
			Position:      "barista",
			BranchId:      MainBranch.Id,
			Salary: &model.Salary{
				Monthly: 3_500_000,
			},
		},
		{
			EmployeeEmail: UserFive.Email,
			Position:      "waiter",
			BranchId:      MainBranch.Id,
			Salary: &model.Salary{
				Monthly: 2_500_000,
			},
		},
	}

	if err := db.CreateInBatches(employees, 5).Error; err != nil {
		log.Fatal(err)
	}

	UserOne.IsOnBoarded = true
	UserTwo.IsOnBoarded = true
	UserThree.IsOnBoarded = true
	UserFour.IsOnBoarded = true
	UserFive.IsOnBoarded = true
	if err := db.Save(&UserOne).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Save(&UserTwo).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Save(&UserThree).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Save(&UserFour).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Save(&UserFive).Error; err != nil {
		log.Fatal(err)
	}

	log.Info("successfully seeded employees")
}
