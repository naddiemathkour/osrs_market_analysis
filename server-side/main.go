package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/naddiemathkour/osrs_market_analysis/db"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize http server to handle requests from client
	srv := startServer()

	// Initialize cron job
	c := startCronJob()

	// Set up signal handling for graceful shut down
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is recieved
	<-s

	// Handle shutdown for cron job
	c.Stop()
	logging.Logger.Info("Scheduler shut down gracefully.")

	// Handle shutdown for HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Fatalf("Http server failed to shut down gracefully: %v", err)
	}

	logging.Logger.Info("Http server shut down gracefully.")
}

func startServer() *http.Server {
	// Initialize Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(db.DataFetchHandler),
	}

	// Create go routine to start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatalf("Http server failed to start: %v", err)
		}
	}()

	logging.Logger.Info("HTTP server started on :8080")
	return srv
}

func startCronJob() *cron.Cron {
	// Create cron object
	c := cron.New()

	// Add function to be run every 5 minutes
	_, err := c.AddFunc("*/5 * * * *", func() {
		logging.Logger.Info("CRON FUNCTION START")
		start := time.Now()
		db.Connect()
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
