package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/test/fixtures"
	"github.com/ilhamSuandi/business_assistant/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestAuthRenewToken(t *testing.T) {
	renewPath := "/api/v1/auth/token/renew"
	loginPath := "/api/v1/auth/login"

	user := model.User{
		Email:    fixtures.UserOne.Email,
		Username: fixtures.UserOne.Username,
		Password: fixtures.UserOne.Password,
	}

	err := helper.CreateUser(&user)
	assert.Nil(t, err)

	t.Run(fmt.Sprintf("GET %s", loginPath), func(t *testing.T) {
		t.Run("should return 200 and renew access token", func(t *testing.T) {
			defer helper.ClearAll(test.DB)

			requestBody := dto.LoginRequest{
				Email:    fixtures.UserOne.Email,
				Password: fixtures.UserOne.Password,
			}

			loginResponse, err := helper.CreateRequest(
				http.MethodPost,
        nil,
				loginPath,
				requestBody,
				nil,
				http.HandlerFunc(authController.Login),
			)
			assert.Nil(t, err)

			refreshToken := loginResponse.Result().Cookies()[0]

			headers := map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"Cookie":       fmt.Sprintf("%s", refreshToken),
			}

			renewResponse, err := helper.CreateRequest(
				http.MethodPost,
        nil,
				renewPath,
				nil,
				&headers,
				http.HandlerFunc(authController.RenewAccessToken),
			)
			assert.Nil(t, err)

			renewResponseBody, err := helper.ParseBody(renewResponse.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, renewResponse.Code)
			assert.Equal(t, "Successfully Renew Access Token", renewResponseBody.Message)
			assert.NotNil(t, renewResponseBody.Data)
		})
	})
}
