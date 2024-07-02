package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	ingest "github.com/naddiemathkour/osrs_market_analysis/db"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create cron job schedular
	c := cron.New()

	// Add function to be run every 5 minutes
	_, err := c.AddFunc("*/5 * * * *", func() {
		logging.Logger.Info("CRON FUNCTION START")
		start := time.Now()
		ingest.Connect()
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

	// Set up signal handling for graceful shut down
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is recieved
	<-s

	// Handle shutdown
	c.Stop()
	logging.Logger.Info("Scheduler shut down gracefully.")
}
