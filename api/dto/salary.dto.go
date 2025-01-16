package dto

type CreateSalaryRequest struct {
	EmployeeEmail string `json:"employee_email" validate:"required"`
	Monthly       *int   `json:"monthly_salary" validate:"omitempty,number"`
	Hourly        *int   `json:"hourly_salary" validate:"omitempty,number"`
}

type CreateSalaryResponse struct {
	EmployeeEmail string `json:"employee_email" example:"userone@gmail.com"`
	Monthly       int    `json:"monthly_salary" example:"3000000"`
	Hourly        int    `json:"hourly_salary" example:"15000"`
}
