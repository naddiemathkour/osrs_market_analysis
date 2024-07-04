package db

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
)

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
