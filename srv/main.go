package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/naddiemathkour/osrs_market_analysis/db"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	serv "github.com/naddiemathkour/osrs_market_analysis/server"
)

func main() {
	// Attempt to connect to database
	d := db.Connect("GET")
	d.Close()

	// Initialize http server to handle requests from client
	srv := serv.StartServer()

	// Initialize cron job
	c := serv.StartCronJob()

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
