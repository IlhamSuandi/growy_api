package dto

type AddEmployeeRequest struct {
	Email       string `json:"email" validate:"required,email" example:"userone@gmail.com"`
	Position    string `json:"position" validate:"required" example:"manager"`
	CompanyName string `json:"company_name" validate:"required" example:"company name"`
	BranchName  string `json:"branch_name" validate:"required" example:"main branch"`
}

type AddEmployeeResponse struct {
	Email    string `json:"email" example:"userone@gmail.com"`
	Position string `json:"position" example:"manager"`
}
