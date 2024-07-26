package db

import (
	"fmt"
	"strings"

	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/naddiemathkour/osrs_market_analysis/models"
)

func Connect(method string) *sqlx.DB {
	// Load env file values
	err := env.Load(".env")
	if err != nil {
		logging.Logger.Fatalf("Error loading env file: %v", err)
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
		logging.Logger.Printf("Failed to connect with: %s", fmt.Sprintf(
			"user=%s dbname=%s sslmode=%s password=%s host=%s port=%s search_path=%s",
			dbConfig.User, dbConfig.Dbname, "disable", dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Path,
		))
		logging.Logger.Errorf("Failed to connect to Postgres: %v", err)
		if strings.Contains(err.Error(), "dial tcp: lookup osrs_db") {
			db.Close()
			return nil
		}
		logging.Logger.Info("Attempting to initialize Database...")
		PostgresInit()
		db.Close()
		return nil
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		logging.Logger.Errorf("Failed to ping database: %v", err)
	}

	return db
}
