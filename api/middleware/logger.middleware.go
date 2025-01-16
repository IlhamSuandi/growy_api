package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/utils"
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

func Log(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := utils.Log
		var path string

		if strings.Contains(r.URL.Path, "/api/v1") {
			path = strings.Split(r.URL.Path, "/api/v1")[1]
		} else {
			path = r.URL.Path
		}

		log.Infof("[%s %s] accepting %s request to %s", r.Method, path, r.Method, r.URL.Path)
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
		userInfo := r.Context().Value("userInfo").(*model.User)

		db.Create(&model.Log{
			Email:         userInfo.Email,
			RemoteAddr:    &ipAddress,
			Action:        &action,
			Method:        &r.Method,
			Path:          &r.URL.Path,
			ExecutionTime: executionTime,
			Status:        uint(responseData.status),
			Size:          responseData.size,
			UserAgent:     r.Header.Get("User-Agent"),
		})

		log.Infof("[%s %s] request finished at %s", r.Method, path, executionTime)
	})
}
