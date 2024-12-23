package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/LitPad/backend/config"
	"github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func RunJobs(cfg config.Config, db *gorm.DB) {
	// RunJobs runs the jobs
	// Initialize the Asynq client and GORM DB (replace with your actual setup)
    redisClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisUrl})

	SetupWorker(db, cfg.RedisUrl)

	// Initial run
	go ReminderJob(db, redisClient)

	// RunWithCron(cfg, db, redisClient)
	RunWithTicker(cfg, db, redisClient)
}

func SetupWorker(db *gorm.DB, redisUrl string) {
	// Set up the Asynq worker
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisUrl},
		asynq.Config{
			Concurrency: 10, // Number of concurrent workers
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// Create a new ServeMux and register task handlers
	mux := asynq.NewServeMux()
	taskHandler := EmailTaskHandler(db)
	mux.HandleFunc(TypeSendEmail, taskHandler)

	// Start the Asynq worker in a separate goroutine to process tasks
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run Asynq server: %v", err)
		}
	}()

}

func RunWithTicker(cfg config.Config, db *gorm.DB, redisClient *asynq.Client) {
	ticker := time.NewTicker(time.Duration(cfg.ReminderCronHours) * time.Hour)

	// Run the job every 2 weeks
	go func() {
		for {
			<-ticker.C
			// Run the ReminderJob function every 2 weeks
			go ReminderJob(db, redisClient)
		}
	}()
}

func RunWithCron(cfg config.Config, db *gorm.DB, redisClient *asynq.Client) {
	// Initialize the cron scheduler
    c := cron.New()

    // Schedule the ReminderJob to run every two weeks at midnight (or any interval)
	cronTime := fmt.Sprintf("@every %sh", cfg.ReminderCronHours)
    c.AddFunc(cronTime, func() {
        go ReminderJob(db, redisClient)
    })

    // Start the cron scheduler
    c.Start()
}