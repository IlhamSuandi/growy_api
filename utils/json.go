package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ParseJSON parses the request body as JSON and stores the result in the value pointed to by v
// It returns an error if the request body is not valid JSON or if it fails to parse the JSON.
func ParseJSON(r *http.Request, v interface{}) error {
	// Validate the request body
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	// Decode the JSON into v
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Validate the request body
	if validationErrors := ValidateStruct(v); validationErrors != nil {
		return validationErrors
	}

	return nil
}
