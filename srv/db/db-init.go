package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofor-little/env"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/naddiemathkour/osrs_market_analysis/models"
)

func PostgresInit() {
	logging.Logger.Infof("Initiate Postgres Database...")
	// Create struct to ingest .env var data
	dbConfig := models.EnvVars{}

	file, err := os.Open(".env")
	if err != nil {
		logging.Logger.Errorf("Error opening env: %v", err)
	}

	env.Load(file.Name())

	// Set values for env file
	dbConfig.Dbname = env.Get("dbname", "")
	dbConfig.Host = env.Get("host", "")
	dbConfig.Path = env.Get("path", "")
	dbConfig.Port = env.Get("port", "")
	dbConfig.User = env.Get("user", "")
	dbConfig.Password = env.Get("password", "")

	// Open connection to database
	logging.Logger.Infof("Attempting to connect to db...")
	logging.Logger.Infof("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.User)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.User)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logging.Logger.Errorf("Error connecting to Postgres: %v\n", err)
		dbClose(db)
	}

	dbQuery := fmt.Sprintf("CREATE DATABASE %s;", dbConfig.Dbname)
	_, err = db.Exec(dbQuery)
	if err != nil {
		logging.Logger.Errorf("Unable to create new database %s: %v", dbConfig.Dbname, err)
	}

	// Connect to new database
	logging.Logger.Infof("Attempting to connect to db...")
	logging.Logger.Infof("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Dbname)
	connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		logging.Logger.Errorf("Error connecting to Postgres: %v\n", err)
		dbClose(db)
	}

	// Check if connection was successful
	err = db.Ping()
	if err != nil {
		logging.Logger.Infof("Error pinging Postgres: %v\n", err)
	} else {
		logging.Logger.Infof("Connection successful!")
	}

	logging.Logger.Infof("Creating database schema...")

	// Create tables and views in Postgres
	tQuery := `CREATE TABLE mapping
					(
						id       int PRIMARY KEY,
						name     varchar(250),
						icon     varchar(250),
						members  boolean,
						lowalch  int,
						highalch int,
						buylimit int,
						value    int,
						examine  varchar(500),
						UNIQUE (id)
					);`

	_, err = db.Exec(tQuery)
	if err != nil {
		logging.Logger.Errorf("Failed to create table: %v\n", err)
		dbClose(db)
	}

	logging.Logger.Infof("Successfully created mapping table")

	tQuery = `CREATE TABLE price
					(
						id              int       NOT NULL,
						timestamp       TIMESTAMP NOT NULL,
						avgHighPrice    int       NOT NULL,
						highPriceVolume int       NOT NULL,
						avgLowPrice     int       NOT NULL,
						lowPriceVolume  int       NOT NULL,
						UNIQUE (id, timestamp)
					);`

	_, err = db.Exec(tQuery)
	if err != nil {
		logging.Logger.Errorf("Failed to create table: %v\n", err)
		dbClose(db)
	}

	logging.Logger.Infof("Successfully created price table")

	vQuery := `CREATE VIEW Listings
				(id, name, icon, examine, members, buylimit, highalch, lowalch, timestamp, avghighprice, highpricevolume, avglowprice, lowpricevolume, spread, margin)
	as
		SELECT DISTINCT m.id,
    m.name,
    m.icon,
    m.examine,
    m.members,
    m.buylimit,
    m.highalch,
    m.lowalch,
    p."timestamp",
        CASE
            WHEN p.avghighprice = 0 THEN p.avglowprice
            WHEN p.avghighprice < p.avglowprice THEN p.avglowprice
            ELSE p.avghighprice
        END AS avghighprice,
        CASE
            WHEN p.avghighprice >= p.avglowprice THEN p.highpricevolume
            ELSE p.lowpricevolume
        END AS highpricevolume,
        CASE
            WHEN p.avglowprice = 0 THEN p.avghighprice
            WHEN p.avglowprice > p.avghighprice THEN p.avghighprice
            ELSE p.avglowprice
        END AS avglowprice,
        CASE
            WHEN p.avglowprice >= p.avghighprice THEN p.highpricevolume
            ELSE p.lowpricevolume
        END AS lowpricevolume,
        CASE
            WHEN p.avghighprice = 0 OR p.avglowprice = 0 THEN 0
            WHEN p.avglowprice > p.avghighprice THEN p.avglowprice - p.avghighprice
            ELSE p.avghighprice - p.avglowprice
        END AS spread,
        CASE
            WHEN p.avghighprice = 0 OR p.avglowprice = 0 THEN 0::double precision
            WHEN p.avglowprice > p.avghighprice THEN round((p.avglowprice - p.avghighprice)::numeric / p.avghighprice::numeric * 100::numeric, 2)::double precision
            ELSE round((p.avghighprice - p.avglowprice)::numeric / p.avglowprice::numeric * 100::numeric, 2)::double precision
        END AS margin
   FROM price p
     JOIN mapping m ON m.id = p.id
     JOIN ( SELECT price.id,
            max(price."timestamp") AS max_timestamp
           FROM price
          GROUP BY price.id) max_ts ON p.id = max_ts.id AND p."timestamp" = max_ts.max_timestamp`

	_, err = db.Exec(vQuery)
	if err != nil {
		logging.Logger.Errorf("Failed to create view: %v\n", err)
		dbClose(db)
	}

	// Get search_path for env
	dbConfig.Path = "public"
	logging.Logger.Infof("Successfully created view!")
	logging.Logger.Infof("All Database requirements executed successfully. Data fetching and storage is now running.")

	db.Close()
	MapItems()
	MapPrices()
}

func checkErr(err error) {
	if err != nil {
		logging.Logger.Infof("Error reading user input: %v\n", err)
		os.Exit(1)
	}
}

func dbClose(db *sql.DB) {
	db.Close()
	os.Exit(-1)
}
