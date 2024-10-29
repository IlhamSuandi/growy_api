package integraion

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	// Define the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	t.Run("GET /_health", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/_health", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(handler)

		handler.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				response.Code, http.StatusOK)
		}

		expected := "OK"
		if response.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				response.Body.String(), expected)
		}
	})
}
