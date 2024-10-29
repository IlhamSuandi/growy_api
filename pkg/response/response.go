package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/types"
)

func WriteJSON(w http.ResponseWriter, status int, response any) error {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Status", fmt.Sprint(status))
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

// WriteError writes an error response to the response writer.
// It sets the status code and writes the error response as JSON.
func WriteError(w http.ResponseWriter, status int, response types.ErrorResponse) {
	WriteJSON(w, status, response)
}
