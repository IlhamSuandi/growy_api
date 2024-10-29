package integration

import (
	"net/http"
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/test/fixtures"
	"github.com/ilhamSuandi/business_assistant/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestAuthLogin(t *testing.T) {
	t.Run("POST /api/v1/auth/login", func(t *testing.T) {
		defer helper.ClearAll(test.DB)

		loginPath := "/api/v1/auth/login"

		requestBody := dto.LoginRequest{
			Email:    fixtures.UserOne.Email,
			Password: fixtures.UserOne.Password,
		}

		user := model.User{
			Email:    fixtures.UserOne.Email,
			Username: fixtures.UserOne.Username,
			Password: fixtures.UserOne.Password,
		}

		err := helper.CreateUser(&user)
		assert.Nil(t, err)

		t.Run("should return 200 and login user", func(t *testing.T) {
			assert.Nil(t, err)
			_, response, err := helper.CreateRequest(http.MethodPost, loginPath, requestBody, authController.Login)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successfully Login", responseBody.Message)
			assert.NotNil(t, responseBody.Data)
		})

		t.Run("should return 401 if user does not exist", func(t *testing.T) {
			requestBody := dto.LoginRequest{
				Email:    fixtures.UserTwo.Email,
				Password: fixtures.UserTwo.Password,
			}

			_, response, err := helper.CreateRequest(http.MethodPost, loginPath, requestBody, authController.Login)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, "email or password is incorrect", responseBody.Message)
		})
	})
}
