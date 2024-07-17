package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gofor-little/env"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	"github.com/naddiemathkour/osrs_market_analysis/models"
)

func PostgresInit() {
	fmt.Println("Initiate Postgres Database...")
	// Create struct to ingest .env var data
	dbConfig := models.EnvVars{}

	// Create console reader:
	reader := bufio.NewReader(os.Stdin)

	// Set values for env file
	dbConfig.Host = "localhost"

	var err error

	fmt.Println("Enter Database user:")
	dbConfig.User, err = reader.ReadString('\n')
	checkErr(err)
	dbConfig.User = strings.TrimSpace(dbConfig.User)

	fmt.Println("Enter Database user password:")
	dbConfig.Password, err = reader.ReadString('\n')
	checkErr(err)
	dbConfig.Password = strings.TrimSpace(dbConfig.Password)

	fmt.Println("Enter port number (default: 5432 or 5433):")
	dbConfig.Port, err = reader.ReadString('\n')
	checkErr(err)
	dbConfig.Port = strings.TrimSpace(dbConfig.Port)

	// Open connection to database
	fmt.Println("Attempting to connect to db...")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error connecting to Postgres: %v\n", err)
		dbClose(db)
	}

	// Check if connection was successful
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging Postgres: %v\n", err)
	} else {
		fmt.Println("Connection successful!")
	}

	// Create new database in Postgres
	fmt.Println("Enter a name for your database:")
	dbConfig.Dbname, err = reader.ReadString('\n')
	checkErr(err)
	dbConfig.Dbname = strings.TrimSpace(dbConfig.Dbname)

	fmt.Printf("Attempting to create database: %s\n", dbConfig.Dbname)

	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbConfig.Dbname)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		fmt.Printf("Failed to create database: %v\n", err)
		dbClose(db)
	}
	fmt.Printf("Database %s created successfully!\n", dbConfig.Dbname)

	// Open connection to database
	fmt.Println("Attempting to connect to new database...")
	connStr = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Port, dbConfig.Dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error connecting to Postgres: %v\n", err)
		dbClose(db)
	}

	// Check if connection was successful
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging Postgres: %v\n", err)
	} else {
		fmt.Println("Connection successful!")
	}

	fmt.Println("Creating database schema...")

	// Create schema in Postgres
	// schemaQuery := fmt.Sprintf("CREATE SCHEMA market_data")
	// _, err = db.Exec(schemaQuery)
	// if err != nil {
	// 	fmt.Printf("Failed to create schema: %v\n", err)
	// 	dbClose(db)
	// }

	// fmt.Println("Schema created successfully!")

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
		fmt.Printf("Failed to create table: %v\n", err)
		dbClose(db)
	}

	fmt.Println("Successfully created mapping table")

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
		fmt.Printf("Failed to create table: %v\n", err)
		dbClose(db)
	}

	fmt.Println("Successfully created price table")

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
		fmt.Printf("Failed to create view: %v\n", err)
		dbClose(db)
	}

	// Get search_path for env
	dbConfig.Path = "public"

	db.Close()
	fmt.Println("Successfully created view!")

	//set values in .env file
	file, err := os.Create(".env")
	if err != nil {
		logging.Logger.Fatalf("Unable to create .env file: %v", err)
	}

	env.Load(file.Name())
	env.Write("host", dbConfig.Host, file.Name(), true)
	env.Write("port", dbConfig.Port, file.Name(), true)
	env.Write("dbname", dbConfig.Dbname, file.Name(), true)
	env.Write("user", dbConfig.User, file.Name(), true)
	env.Write("password", dbConfig.Password, file.Name(), true)
	env.Write("path", dbConfig.Path, file.Name(), true)

	fmt.Println("All Database requirements executed successfully. Data fetching and storage is now running.")
	MapItems()
	MapPrices()
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error reading user input: %v\n", err)
		os.Exit(1)
	}
}

func dbClose(db *sql.DB) {
	db.Close()
	os.Exit(-1)
}
