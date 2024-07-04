package db

import (
	"fmt"

	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
)

// Create environment variable struct
type EnvVars struct {
	host     string
	port     string
	dbname   string
	user     string
	password string
	path     string
}

func Connect() {
	// Log when function begins running
	logging.Logger.Info("Connecting to Database...")

	//load env file from directory
	if err := env.Load("./.env"); err != nil {
		logging.Logger.Fatalf("Failed to load environment file: %v", err)
	}

	dbConfig := EnvVars{
		host:     env.Get("host", ""),
		port:     env.Get("port", ""),
		dbname:   env.Get("dbname", ""),
		user:     env.Get("user", ""),
		password: env.Get("password", ""),
		path:     env.Get("path", ""),
	}

	//connect to PostgreSQL database
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s password=%s host=%s port=%s search_path=%s",
		dbConfig.user, dbConfig.dbname, "disable", dbConfig.password, dbConfig.host, dbConfig.port, dbConfig.path,
	))
	if err != nil {
		logging.Logger.Fatalf("Failed to connect to Postgres: %v", err)
	} else {
		logging.Logger.Info("Successfully Connected.")
	}

	defer db.Close()

	MapItems(db)
	MapPrices(db)
}
