package main

import (
	"fmt"
	"log"

	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	database "github.com/naddiemathkour/osrs_market_analysis/db"
)

type EnvVars struct {
	host     string
	port     string
	dbname   string
	user     string
	password string
}

func main() {
	//load env file from directory
	if err := env.Load("./.env"); err != nil {
		log.Fatal(err)
	}

	dbConfig := EnvVars{
		host:     env.Get("host", ""),
		port:     env.Get("port", ""),
		dbname:   env.Get("dbname", ""),
		user:     env.Get("user", ""),
		password: env.Get("password", ""),
	}

	//connect to PostgreSQL database
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		dbConfig.user, dbConfig.dbname, "disable", dbConfig.password, dbConfig.host, dbConfig.port,
	))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//test connection to database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	database.MapItems(db)
}
