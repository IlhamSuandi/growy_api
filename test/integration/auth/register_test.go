package integration

import (
	"net/http"
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/test/fixtures"
	"github.com/ilhamSuandi/business_assistant/test/helper"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthRegister(t *testing.T) {
	t.Run("POST /api/v1/auth/register", func(t *testing.T) {
		registerPath := "/api/v1/auth/register"
		requestBody := dto.RegisterRequest{
			Email:    fixtures.UserOne.Email,
			Username: fixtures.UserOne.Username,
			Password: fixtures.UserOne.Password,
		}

		t.Run("should return 200 and register user", func(t *testing.T) {
			defer helper.ClearAll(test.DB)
			_, response, err := helper.CreateRequest(http.MethodPost, registerPath, requestBody, authController.Register)
			assert.Nil(t, err)

			var responseData dto.RegisteredUserResponse

			responseBody, err := helper.ParseBody(response.Body, &responseData)
			assert.Nil(t, err)

			expectedResponse := types.Response{
				Message: "Successfully Register User",
				Data: dto.RegisteredUserResponse{
					Username: requestBody.Username,
					Email:    requestBody.Email,
				},
				Status: http.StatusCreated,
			}

			actualResponse := types.Response{
				Message: "Successfully Register User",
				Data: dto.RegisteredUserResponse{
					Username: responseData.Username,
					Email:    responseData.Email,
				},
				Status: http.StatusCreated,
			}

			assert.Equal(t, http.StatusCreated, response.Code)
			assert.Equal(t, "Successfully Register User", responseBody.Message)
			assert.Equal(t, expectedResponse, actualResponse)
		})

		t.Run("should return 400 if email is invalid", func(t *testing.T) {
			defer helper.ClearAll(test.DB)
			requestBody.Email = "invalid"
			_, response, err := helper.CreateRequest(http.MethodPost, registerPath, requestBody, authController.Register)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, response.Code)
		})

		t.Run("should return 400 if password is not contains special char, at least one number and one uppercase letter", func(t *testing.T) {
			defer helper.ClearAll(test.DB)
			requestBody.Email = "ilham@gmail.com"
			requestBody.Password = "invalid"
			_, response, err := helper.CreateRequest(http.MethodPost, registerPath, requestBody, authController.Register)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, response.Code)
		})

		t.Run("should return 409 if user is exists", func(t *testing.T) {
			defer helper.ClearAll(test.DB)
			requestBody.Email = "ilham@gmail.com"
			requestBody.Password = "Ilham123;"
			user := model.User{
				Email:    requestBody.Email,
				Username: requestBody.Username,
			}

			err := helper.CreateUser(&user)
			assert.Nil(t, err)

			_, response, err := helper.CreateRequest(http.MethodPost, registerPath, requestBody, authController.Register)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusConflict, response.Code)
			assert.Equal(t, "Error Registering User", responseBody.Message)
		})
	})
}
