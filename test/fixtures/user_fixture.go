package fixtures

import (
	"github.com/ilhamSuandi/business_assistant/database/model"
)

var UserOne = model.User{
	Username: "testing1",
	Email:    "testing1@gmail.com",
	Password: "Testing1+",
}

var UserTwo = model.User{
	Username: "testing2",
	Email:    "testing2@gmail.com",
	Password: "Testing2+",
}
