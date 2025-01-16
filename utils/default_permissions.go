package utils

import (
	permission "github.com/ilhamSuandi/business_assistant/constant"
	"github.com/ilhamSuandi/business_assistant/database/model"
)

func EmployeeDefaultPermissions() []*model.Permission {
	return []*model.Permission{
		{Resource: permission.Me, Action: "all"},
		{Resource: permission.Qrcode, Action: "get"},
		{Resource: permission.EmployeeRoute, Action: "get"},
	}
}

func OwnerDefaultPermissions() []*model.Permission {
	return []*model.Permission{
		{Resource: permission.Owner, Action: "all"},
	}
}

func AdminDefaultPermissions() []*model.Permission {
	return []*model.Permission{
		{Resource: permission.Admin, Action: "all"},
	}
}
