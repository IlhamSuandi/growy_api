package seeds

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/auth"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

var UserOne = model.User{
	Username:    "user1",
	Email:       "user1@gmail.com",
	Password:    "User123+",
	Role:        "employee",
	Permissions: utils.EmployeeDefaultPermissions(),
}

var UserTwo = model.User{
	Username: "user2",
	Email:    "user2@gmail.com",
	Password: "User123+",
	Role:     "employee",
	Permissions: utils.EmployeeDefaultPermissions(),
}

var UserThree = model.User{
	Username: "user3",
	Email:    "user3@gmail.com",
	Password: "User123+",
	Role:     "employee",
	Permissions: utils.EmployeeDefaultPermissions(),
}

var UserFour = model.User{
	Username: "user4",
	Email:    "user4@gmail.com",
	Password: "User123+",
	Role:     "employee",
	Permissions: utils.EmployeeDefaultPermissions(),
}

var UserFive = model.User{
	Username: "user5",
	Email:    "user5@gmail.com",
	Password: "User123+",
	Role:     "employee",
	Permissions: utils.EmployeeDefaultPermissions(),
}

var UserOwner = model.User{
	Username: "owner",
	Email:    "owner@gmail.com",
	Password: "Owner123+",
	Role:     "owner",
}

var UserAdmin = model.User{
	Username: "admin",
	Email:    "admin@example.com",
	Password: "Securepassword+",
	Role:     "admin",
}

func createUser(db *gorm.DB, user *model.User) error {
	userrepo := repository.NewUserRepository(db)
	hashedPass, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPass)
	return userrepo.CreateUser(user)
}

func SeedUsers(db *gorm.DB) {
	log := utils.Log
	log.Info("seeding users")

	users := []*model.User{
		&UserAdmin,
		&UserOwner,
		&UserOne,
		&UserTwo,
		&UserThree,
		&UserFour,
		&UserFive,
	}

	for _, user := range users {
		log.Infof("creating user %s", user.Email)
		if err := createUser(db, user); err != nil {
			log.Errorf("error creating user %s", user.Email)
			log.Fatal(err)
		}
	}

	log.Info("successfully seeded users")
}
