package serv

import (
	"time"

	db "github.com/naddiemathkour/osrs_market_analysis/db"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func StartCronJob() *cron.Cron {
	// Create cron object
	c := cron.New()

	// Add function to be run every 5 minutes
	_, err := c.AddFunc("*/5 * * * *", func() {
		logging.Logger.Info("CRON FUNCTION START")
		start := time.Now()
		db.Connect("GET")
		logging.Logger.WithFields(logrus.Fields{
			"duration": time.Since(start).Seconds(),
		}).Info("Cron function completed.")
	})
	if err != nil {
		logging.Logger.Fatalf("Error scheduling task: %v", err)
	}

	// Start the cron scheduler
	c.Start()
	logging.Logger.Info("Cron schedular started")
	return c
}
