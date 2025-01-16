package middleware

import (
	"net/http"

	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/utils"
	"golang.org/x/time/rate"
)

func Limiter(limit rate.Limit, burst *int, next http.Handler) http.Handler {
	log := utils.Log

	defaultBurst := int(limit * 2)

	if burst == nil {
		burst = &defaultBurst
	}

	limiter := rate.NewLimiter(limit, *burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			log.Error("[middleware] too many requests")
			response.WriteError(w, http.StatusTooManyRequests, types.ErrorResponse{
				Message: "Too many requests",
				Error:   "Too many requests",
				Status:  http.StatusTooManyRequests,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
