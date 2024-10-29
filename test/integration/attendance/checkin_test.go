package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/api/middleware"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/test/fixtures"
	"github.com/ilhamSuandi/business_assistant/test/helper"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckin(t *testing.T) {
	defer helper.ClearAll(test.DB)
	user := model.User{
		Username: fixtures.UserOne.Username,
		Email:    fixtures.UserOne.Email,
		Password: fixtures.UserOne.Password,
	}

	err := helper.CreateUser(&user)
	assert.Nil(t, err)

	qrData, err := fixtures.GetQrCode(user.UUID)
	assert.Nil(t, err)

	checkinPath := fmt.Sprintf("/api/v1/attendances/check-in/%v", qrData.Code)

	payload := dto.CheckInRequest{
		Location: "Jakarta",
	}

	accessToken, err := fixtures.AccessToken(user.UUID, user.Username, user.Email)
	assert.Nil(t, err)

	checkinHandler := middleware.Auth(http.HandlerFunc(attendanceController.CheckIn), test.DB)

	t.Run("POST /api/v1/attendances/check-in", func(t *testing.T) {
		t.Run("should return 200 and checkin info", func(t *testing.T) {
			requestBody, err := json.Marshal(payload)
			assert.Nil(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				checkinPath,
				strings.NewReader(string(requestBody)),
			)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+accessToken)

			response := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.Handle("/api/v1/attendances/check-in/{token}", checkinHandler)
			mux.ServeHTTP(response, request)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successfully Checked In Attendance", responseBody.Message)
		})

		t.Run("should return 403 if token is empty", func(t *testing.T) {
			requestBody, err := json.Marshal(payload)
			assert.Nil(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/attendances/check-in",
				strings.NewReader(string(requestBody)),
			)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+accessToken)

			response := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.Handle("/api/v1/attendances/check-in", checkinHandler)
			mux.ServeHTTP(response, request)

			responseBody, err := helper.ParseBody(response.Body, &types.Response{})
			assert.Nil(t, err)

			assert.Equal(t, http.StatusForbidden, response.Code)
			assert.Equal(t, "token is required", responseBody.Message)
		})

		t.Run("should return 401 if token is not valid", func(t *testing.T) {
			requestBody, err := json.Marshal(payload)
			assert.Nil(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				checkinPath+"invalid-token",
				strings.NewReader(string(requestBody)),
			)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+accessToken)

			response := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.Handle("/api/v1/attendances/check-in/{token}", checkinHandler)
			mux.ServeHTTP(response, request)

			responseBody, err := helper.ParseBody(response.Body, &types.Response{})
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, "error checking in attendance", responseBody.Message)
		})

		t.Run("should return 401 if Authorization header is not valid", func(t *testing.T) {
			requestBody, err := json.Marshal(payload)
			assert.Nil(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				checkinPath,
				strings.NewReader(string(requestBody)),
			)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			response := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.Handle("/api/v1/attendances/check-in/{token}", checkinHandler)
			mux.ServeHTTP(response, request)

			responseBody, err := helper.ParseBody(response.Body, &types.Response{})
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, "Unauthorized", responseBody.Message)
		})

		t.Run("should return 400 if request body is empty", func(t *testing.T) {
			accessToken, err := fixtures.AccessToken(user.UUID, user.Username, user.Email)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, checkinPath, nil)
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+accessToken)

			response := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.Handle("/api/v1/attendances/check-in/{token}", checkinHandler)
			mux.ServeHTTP(response, request)

			responseBody, err := helper.ParseBody(response.Body, nil)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "error parsing json", responseBody.Message)
		})
	})
}
