package main

import (
	"github.com/ilhamSuandi/business_assistant/cmd/seeder/seeds"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/database"
	"github.com/ilhamSuandi/business_assistant/utils"
)

func main() {
	log := utils.Log
	log.Info("===Seeder Staerted===")
	db, err := database.Connect(config.DB_HOST, config.DB_NAME)
	database.AutoMigrate(db)

	if err != nil {
		panic("failed to connect to database")
	}

	seeds.SeedUsers(db)
	seeds.SeedQrCodes(db)
	seeds.SeedCompany(db)
	seeds.SeedBranch(db)
	seeds.SeedEmployees(db)
	// seeds.SeedRole(db)
	seeds.SeedCheckins(db)

	log.Info("===Seeder Finished===")
}
