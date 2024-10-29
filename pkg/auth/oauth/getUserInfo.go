package oauth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/types"
)

func GetUserInfo(accessToken string) (*types.GoogleUser, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo?access_token="+accessToken, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userData := types.GoogleUser{}

	if err := json.Unmarshal(respBody, &userData); err != nil {
		return nil, err
	}

	return &userData, nil
}
