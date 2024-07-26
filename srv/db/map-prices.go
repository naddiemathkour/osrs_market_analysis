package db

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	models "github.com/naddiemathkour/osrs_market_analysis/models"
)

func MapPrices() {
	// Get db connection
	db := Connect("POST")
	if db == nil {
		return
	}
	// Log operation
	logging.Logger.Info("Ingesting item price data...")

	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/5m", nil)
	if err != nil {
		logging.Logger.Errorf("Failed to create http request: %v", err)
		db.Close()
		return
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.Logger.Errorf("Failed to accept request: %v", err)
		db.Close()
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Logger.Errorf("Failed to read response body: %v", err)
		db.Close()
		return
	}

	//Decode response data to JSON
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logging.Logger.Errorf("Failed to unmarshal JSON data: %v", err)
		db.Close()
		return
	}

	jsonResp, err := json.Marshal(data["data"])
	if err != nil {
		logging.Logger.Errorf("Failed to marshal data: %v", err)
		db.Close()
		return
	}

	// Unmarshal data into item map struct map[string]models.ItemData
	var itemsMap map[string]models.ItemData
	err = json.Unmarshal(jsonResp, &itemsMap)
	if err != nil {
		logging.Logger.Errorf("Failed to unmarshal data: %v", err)
		db.Close()
		return
	}

	// Convert map into slice of Items
	var items []models.ItemObject
	for id, data := range itemsMap {
		items = append(items, models.ItemObject{
			ID:   id,
			Data: data,
		})
	}

	// Set timestamp for ingestion
	var timestamp = time.Now().UTC().Format("2006-01-02 15:04:05")

	// Iterate through items and insert into database
	count := 0
	for _, item := range items {
		insertQuery := `INSERT INTO public.price (id, timestamp, avgHighPrice, highPriceVolume, avgLowPrice, lowPriceVolume)
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (id, timestamp) DO NOTHING;`

		_, err := db.Exec(insertQuery, item.ID, timestamp, item.Data.AvgHighPrice, item.Data.HighPriceVolume, item.Data.AvgLowPrice, item.Data.LowPriceVolume)
		if err != nil {
			logging.Logger.Errorf("Failed to execute insert query: %v", err)
			db.Close()
			return
		} else {
			count++
		}
	}

	deleteQuery := `DELETE FROM price WHERE timestamp < NOW() -'15 minutes'::interval;`
	_, err = db.Exec(deleteQuery)
	if err != nil {
		logging.Logger.Errorf("Failed to delete stale data: %v", err)
		db.Close()
		return
	}

	db.Close()

	logging.Logger.Infof("Successfully inserted %d items price data.", count)
}
