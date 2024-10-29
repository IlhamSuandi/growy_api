package jobs

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ilhamSuandi/business_assistant/database/model"
	qr "github.com/ilhamSuandi/business_assistant/pkg/qrcode"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

type Scheduler struct {
	cron gocron.Scheduler
}

var Jobs gocron.Scheduler

func generateQrCodes(db *gorm.DB, workerCount int) {
	userRepo := repository.NewUserRepository(db)

	// Get all users
	users, err := userRepo.GetUsers(nil, nil)
	if err != nil {
		utils.Log.Panic(err)
	}

	wg := sync.WaitGroup{}
	workerWg := sync.WaitGroup{}
	userChan := make(chan model.User)

	// start worker
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go workerTask(db, &workerWg, userChan)
	}

	wg.Add(len(users))

	go func() {
		for _, user := range users {
			userChan <- user
		}
		close(userChan)
	}()

	wg.Wait()
	workerWg.Wait()
}

func workerTask(db *gorm.DB, workerWg *sync.WaitGroup, userChan chan model.User) {
	defer workerWg.Done()

	for user := range userChan {
		// Start a new transaction for each user
		tx := db.Begin()
		if tx.Error != nil {
			utils.Log.Error(tx.Error)
			return
		}

		qrcode := model.QRCode{}
		tx.Where("UUID = ?", user.UUID).First(&qrcode)

		// Generate QR code and expiration date
		code := qr.GenerateRandomCode(user.UUID)
		expiredAt := qr.GetExpirationDate()

		// Update QRCode struct
		qrcode.Code = code
		qrcode.ExpiresAt = &expiredAt
		qrcode.IsUsed = false

		// Save the updated QR code to the database
		if err := tx.Save(&qrcode).Error; err != nil {
			utils.Log.Error(err)
			tx.Rollback()
			continue
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			utils.Log.Error(err)
			continue
		}

		utils.Log.Infof("successfully regenerated qr code for user %s", user.Email)
	}
}

func CreateScheduler(location *time.Location) {
	cron, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		utils.Log.Panic("Error creating scheduler")
	}

	Jobs = cron
}

func Start(db *gorm.DB) {
	utils.Log.Info("starting scheduler")
	_, err := Jobs.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(6, 30, 0),
			),
		),
		gocron.NewTask(
			func(a string, b int) {
				utils.Log.Info("running scheduler")
				generateQrCodes(db, 100)
				defer utils.Log.Info("QR Code Generated")
			},
			"qr-code-generated",
			1,
		),
		gocron.WithName("Qr Code Handler"),
		gocron.WithTags("Absent Qr Code"),
	)
	if err != nil {
		utils.Log.Panic(err)
	}

	Jobs.Start()
}

func Stop() {
	if err := Jobs.Shutdown(); err != nil {
		utils.Log.Errorf("Error shutting down scheduler: %v", err)
	}
	utils.Log.Info("scheduler stopped")
}

func GetAllJobs() []gocron.Job {
	return Jobs.Jobs()
}
