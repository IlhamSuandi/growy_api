package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/utils"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := utils.Log
		// get allowed origins from env
		log.Info("getting allowed origins from env")
		allowedOrigins := config.ALLOWED_ORIGIN

		// change string to env
		origins := strings.Split(allowedOrigins, ",")

		// get origin from request
		log.Info("getting origin from request")
		reqOrigin := r.Header.Get("Origin")

		// allow all localhost
		log.Info("allowing all localhost")
		if strings.HasPrefix(reqOrigin, "http://localhost") || strings.HasPrefix(reqOrigin, "http://127.0.0.1") {
			origins = append(origins, reqOrigin)
		}

		log.Info("checking if origin is allowed")
		if slices.Contains(origins, reqOrigin) {
			w.Header().Set("Access-Control-Allow-Origin", reqOrigin)
		} else {
			// if origin not allowed, return forbidden
			log.Error("origin not allowed")
			response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
				Error:   "Origin not allowed",
				Message: "Origin not allowed",
				Status:  http.StatusForbidden,
			})
			return
		}

		// Allow specific HTTP methods
		log.Info("setting headers")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

		// if origin is allowed, continue
		next.ServeHTTP(w, r)
	})
}
