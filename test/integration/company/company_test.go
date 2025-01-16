package integration

import (
	"net/http"
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/api/middleware"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/test/fixtures"
	"github.com/ilhamSuandi/business_assistant/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestCompany(t *testing.T) {
	defer helper.ClearAll(test.DB)
	user := model.User{
		Username: fixtures.UserOne.Username,
		Email:    fixtures.UserOne.Email,
		Password: fixtures.UserOne.Password,
	}

	err := helper.CreateUser(&user)
	assert.Nil(t, err)

	companyPath := "/api/v1/company"

	accessToken, err := fixtures.AccessToken(user.UUID, user.Username, user.Email)
	assert.Nil(t, err)

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	companyHandler := middleware.Auth(http.HandlerFunc(CompanyController.CreateCompany), test.DB)

	t.Run("POST /api/v1/company", func(t *testing.T) {
		t.Run("should return 200 and create company", func(t *testing.T) {
			payload := dto.CreateCompanyRequest{
				Name:    "Growy",
				Address: "Jakarta",
			}

			response, err := helper.CreateRequest(
				http.MethodPost,
				nil,
				companyPath,
				payload,
				&headers,
				companyHandler,
			)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successfully created company", responseBody.Message)
		})

		t.Run("should return 401 unauthorized", func(t *testing.T) {
			payload := dto.CreateCompanyRequest{
				Name:    "Growy",
				Address: "Jakarta",
			}

			response, err := helper.CreateRequest(
				http.MethodPost,
				nil,
				companyPath,
				payload,
				nil,
				companyHandler,
			)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, "Unauthorized", responseBody.Message)
		})

		t.Run("should return 400 if name is empty", func(t *testing.T) {
			payload := dto.CreateCompanyRequest{
				Address: "Jakarta",
			}

			response, err := helper.CreateRequest(
				http.MethodPost,
				nil,
				companyPath,
				payload,
				&headers,
				companyHandler,
			)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "Error request parsing body", responseBody.Message)
		})

		t.Run("should return 400 if address is empty", func(t *testing.T) {
			payload := dto.CreateCompanyRequest{
				Name: "Growy",
			}

			response, err := helper.CreateRequest(
				http.MethodPost,
				nil,
				companyPath,
				payload,
				&headers,
				companyHandler,
			)
			assert.Nil(t, err)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "Error request parsing body", responseBody.Message)
		})
	})
}
