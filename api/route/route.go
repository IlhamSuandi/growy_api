package route

import (
	"net/http"

	"github.com/ilhamSuandi/business_assistant/api/middleware"
	"github.com/ilhamSuandi/business_assistant/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) *http.ServeMux {
	// Router Config
	router := http.NewServeMux()

	v1 := http.NewServeMux()
	protected := http.NewServeMux()

	// grouping routes to "/api/v1"
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// use auth middleware
	v1.Handle("/", middleware.Auth(protected, db))

	// Health Check
	router.HandleFunc("GET /_health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		},
	)

	if config.APP_ENV == "development" {
		swagger := http.NewServeMux()

		// grouping routes to "/swagger"
		router.Handle("/swagger/", http.StripPrefix("/swagger", swagger))

		router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./docs/swagger.json")
		})
		swagger.HandleFunc("/",
			httpSwagger.WrapHandler,
		)

	}

	// Public Routes
	AuthRoutes(v1, db, "/auth")

	// Protected Routes
	UserRoutes(protected, db, "/users")
	AttendanceRoutes(protected, db, "/attendances")

	// TODO: create role middleware
	QRCodeRoutes(protected, db, "/qrcode")

	return router
}
