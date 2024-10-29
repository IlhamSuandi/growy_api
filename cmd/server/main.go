package main

import (
	"time"

	"github.com/ilhamSuandi/business_assistant/api"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/database"
	"github.com/ilhamSuandi/business_assistant/pkg/jobs"
	"github.com/ilhamSuandi/business_assistant/utils"
)

// @title Growy API
// @version 1.0.0
// @description Growy API

// @host localhost:5000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Example Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
// query.collection.format multi

// @contact.name Ilham Suandi
// @contact.url https://github.com/ilhamSuandi
// @contact.email ilham15suandi@gmail.com
func main() {
	// Connect to DB
	log := utils.Log
	db, err := database.Connect(config.DB_HOST, config.DB_NAME)
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	// Close DB Connection at the end of the program
	defer database.CloseDb(db)

	// AutoMigrate
	database.AutoMigrate(db)

	// Create Cron Job
	location, _ := time.LoadLocation("Asia/Jakarta")
	jobs.CreateScheduler(location)
	jobs.Start(db)

	// Stop Cron Job at the end of the program
	defer jobs.Stop()

	utils.RegisterValidator()
	server := api.NewApiServer(":5000", db)

	// Start the server in a goroutine
	server.Start()

	// Set up graceful shutdown
	server.GracefulShutdown()
}
