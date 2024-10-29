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

	if err := seeds.SeedUsers(db); err != nil {
		log.Fatal(err)
	}
	if err := seeds.SeedQrCodes(db); err != nil {
		log.Fatal(err)
	}
	if err := seeds.SeedCompany(db); err != nil {
		log.Fatal(err)
	}
	if err := seeds.SeedBranch(db); err != nil {
		log.Fatal(err)
	}
	if err := seeds.SeedRole(db); err != nil {
		log.Fatal(err)
	}
	if err := seeds.SeedCheckins(db); err != nil {
		log.Fatal(err)
	}

	log.Info("===Seeder Finished===")
}
