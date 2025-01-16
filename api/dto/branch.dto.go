package dto

type CreateBranchRequest struct {
	CompanyName string `json:"company_name" validate:"required" example:"company name"`
	BranchName  string `json:"branch_name" validate:"required" example:"branch name"`
	Address     string `json:"address" validate:"required" example:"jakarta"`
}
