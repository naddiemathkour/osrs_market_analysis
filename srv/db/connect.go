package db

import (
	"fmt"

	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/naddiemathkour/osrs_market_analysis/models"
)

func Connect(method string) *sqlx.DB {
	//load env file from directory
	if err := env.Load("./.env"); err != nil {
		logging.Logger.Info("Failed to load environment file. Initiating DB setup...")
		PostgresInit()
	}

	dbConfig := models.EnvVars{
		Host:     env.Get("host", ""),
		Port:     env.Get("port", ""),
		Dbname:   env.Get("dbname", ""),
		User:     env.Get("user", ""),
		Password: env.Get("password", ""),
		Path:     env.Get("path", ""),
	}

	//connect to PostgreSQL database
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s password=%s host=%s port=%s search_path=%s",
		dbConfig.User, dbConfig.Dbname, "disable", dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Path,
	))
	if err != nil {
		logging.Logger.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		logging.Logger.Errorf("Failed to ping database: %v", err)
	}

	return db
}
