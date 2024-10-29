package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/database/model"
)

type CheckInRequest struct {
	Location string `json:"location" validate:"required,min=3" example:"Jakarta"`
}

type CheckInResponse struct {
	UserId   uuid.UUID  `json:"user_id" example:"00000000-0000-0000-0000-000000000000"`
	Status   string     `json:"status" example:"present"`
	CheckIn  *time.Time `json:"check_in" example:"2022-01-01T00:00:00Z"`
	Date     string     `json:"date" example:"2022-01-01"`
	Location *string    `json:"location" example:"Jakarta"`
}

type CheckoutRequest struct {
	UserId   uuid.UUID  `json:"user_id"`
	Status   string     `json:"status"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
	Date     string     `json:"date"`
	Location *string    `json:"location"`
}

type GetAttendancesResponse struct {
	UserId      uuid.UUID           `json:"user_id"`
	Attendances []*model.Attendance `json:"attendances"`
}
