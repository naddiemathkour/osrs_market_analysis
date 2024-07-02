package ingest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofor-little/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
)

// Create mapping object struct with exported fields and JSON tags
type MapObject struct {
	ID       interface{} `json:"id"`
	Name     interface{} `json:"name"`
	Members  interface{} `json:"members"`
	Highalch interface{} `json:"highalch"`
	Lowalch  interface{} `json:"lowalch"`
	Value    interface{} `json:"value"`
	Buylimit interface{} `json:"limit"`
	Icon     interface{} `json:"icon"`
	Examine  interface{} `json:"examine"`
}

// Create item price and data structs
type ItemData struct {
	AvgHighPrice    int `json:"avgHighPrice"`
	AvgLowPrice     int `json:"avgLowPrice"`
	HighPriceVolume int `json:"highPriceVolume"`
	LowPriceVolume  int `json:"lowPriceVolume"`
}

type Item struct {
	ID   string   `json:"id"`
	Data ItemData `json:"data"`
}

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

func MapItems(db *sqlx.DB) {
	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/mapping", nil)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.Logger.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	//Decode response data to JSON
	var jsonResp []map[string]interface{}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	// Query database for item count. If item count is still the same, return. Else, update all item data.
	countQuery := `SELECT COUNT(*) OUTPUT FROM mapping`
	var count int
	err = db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		logging.Logger.Fatalf("Error querying row count: %v", err)
	}

	if len(jsonResp) == count {
		return
	}

	// Log when attempting to add items
	logging.Logger.Infof("Adding %d item(s) to Database...", len(jsonResp)-count)

	// Upsert all item data
	for _, item := range jsonResp {
		// Gather required data for insertion
		tempObj := MapObject{
			ID:       item["id"],
			Name:     item["name"],
			Members:  item["members"],
			Highalch: item["highalch"],
			Lowalch:  item["lowalch"],
			Value:    item["value"],
			Buylimit: item["limit"],
			Icon:     item["icon"],
			Examine:  item["examine"],
		}

		// Create and Execute INSERT statement
		insertQuery := `INSERT INTO mapping (id, members, lowalch, highalch, buylimit, value, icon, name, examine)
                    	VALUES (:id, :members, :lowalch, :highalch, :buylimit, :value, :icon, :name, :examine)
						ON CONFLICT (id)
						DO UPDATE SET
							members 	= EXCLUDED.members,
							lowalch 	= EXCLUDED.lowalch,
							highalch 	= EXCLUDED.highalch,
							buylimit 	= EXCLUDED.buylimit,
							value 		= EXCLUDED.value,
							icon 		= EXCLUDED.icon,
							name 		= EXCLUDED.name,
							examine 	= EXCLUDED.examine;`

		_, err := db.NamedExec(insertQuery, &tempObj)
		if err != nil {
			logging.Logger.Fatal(err)
		}
	}

	logging.Logger.Info("Successfully added item updates.")
}

func MapPrices(db *sqlx.DB) {
	// Log operation
	logging.Logger.Info("Ingesting item price data...")

	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/5m", nil)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.Logger.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	//Decode response data to JSON
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	jsonResp, err := json.Marshal(data["data"])
	if err != nil {
		logging.Logger.Fatal(err)
	}

	// Unmarshal data into item map struct map[string]ItemData
	var itemsMap map[string]ItemData
	err = json.Unmarshal(jsonResp, &itemsMap)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	// Convert map into slice of Items
	var items []Item
	for id, data := range itemsMap {
		items = append(items, Item{
			ID:   id,
			Data: data,
		})
	}

	// Set timestamp for ingestion
	var timestamp = time.Now().Format("2006-01-02 15:04:05")

	// Iterate through items and insert into database
	count := 0
	for _, item := range items {
		insertQuery := `INSERT INTO price (id, timestamp, avgHighPrice, highPriceVolume, avgLowPrice, lowPriceVolume)
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (id, timestamp) DO NOTHING;`

		_, err := db.Exec(insertQuery, item.ID, timestamp, item.Data.AvgHighPrice, item.Data.HighPriceVolume, item.Data.AvgLowPrice, item.Data.LowPriceVolume)
		if err != nil {
			logging.Logger.Fatal(err)
		} else {
			count++
		}
	}

	logging.Logger.Infof("Successfully inserted %d items price data.", count)
}
