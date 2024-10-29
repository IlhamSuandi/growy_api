package middleware

import (
	"net/http"
	"time"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"gorm.io/gorm"
)

type (
	responseData struct {
		status int
		size   int
	}

	Logger struct {
		time          time.Time
		remoteAddr    string
		method        string
		path          string
		executionTime time.Duration
		status        int
		size          int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Log(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		next.ServeHTTP(&lrw, r)

		executionTime := time.Since(start).String()

		var ipAddress string

		if ipAddress = r.Header.Get("X-Forwarded-For"); ipAddress == "" {
			ipAddress = r.RemoteAddr
		}

		action := "request"

		db.Create(&model.Log{
			RemoteAddr:    &ipAddress,
			Action:        &action,
			Method:        &r.Method,
			Path:          &r.URL.Path,
			ExecutionTime: executionTime,
			Status:        uint(responseData.status),
			Size:          responseData.size,
			UserAgent:     r.Header.Get("User-Agent"),
		})
	})
}
