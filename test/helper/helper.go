package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/test"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateRequest(method string,
	path string,
	requestBody interface{},
	handler func(http.ResponseWriter, *http.Request),
) (*http.Request, *httptest.ResponseRecorder, error) {
	bodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, err
	}

	request := httptest.NewRequest(http.MethodPost,
		path,
		strings.NewReader(string(bodyJson)),
	)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response := httptest.NewRecorder()

	http.HandlerFunc(handler).ServeHTTP(response, request)

	return request, response, err
}

func ParseBody(responseBody io.Reader, targetType interface{}) (*types.Response, error) {
	// Read all bytes from the response body
	bytes, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	// Create a Response struct to hold the parsed data
	body := types.Response{}

	// Unmarshal the entire response body into the Response struct
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return nil, err
	}

	// Marshal the Data field of the Response into JSON bytes
	dataBytes, err := json.Marshal(body.Data)
	if err != nil {
		return nil, err
	}

	// Unmarshal the Data JSON bytes into the specified target type
	if err := json.Unmarshal(dataBytes, &targetType); err != nil {
		return nil, err
	}

	// Set the target type as the Data in the response
	body.Data = targetType

	// Return the populated Response struct
	return &body, nil
}

func ParseData(targetType interface{}, responseBody interface{}) error {
	responseDataBytes, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(responseDataBytes, &targetType); err != nil {
		return err
	}

	return nil
}

func CreateUser(user *model.User) error {
	userRepo := repository.NewUserRepository(test.DB)
	sessionRepo := repository.NewSessionRepository(test.DB)

	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo)

	if _, isUserExists := authUsecase.IsUserExists(user.Email); isUserExists {
		return errors.New("user already exists")
	}

	if err := authUsecase.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func ClearAll(db *gorm.DB) {
	tables := []interface{}{
		model.User{},
		model.Attendance{},
		model.QRCode{},
	}

	for _, table := range tables {
		ClearData(db, table)
	}
}

func ClearData(db *gorm.DB, model interface{}) {
	err := db.Where("id is not null").Delete(&model).Error
	if err != nil {
		logrus.Fatalf("Failed clear data : %+v", err)
	}
}
