package db

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
	models "github.com/naddiemathkour/osrs_market_analysis/models"
)

func MapItems() {
	// Get db connection
	db := Connect("POST")

	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/mapping", nil)
	if err != nil {
		logging.Logger.Fatalf("Failed to create http request: %v", err)
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.Logger.Fatalf("Failed to accept request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Logger.Fatalf("Failed to read response body: %v", err)
	}

	//Decode response data to JSON
	var jsonResp []map[string]interface{}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		logging.Logger.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Query database for item count. If item count is still the same, return. Else, update all item data.
	countQuery := `SELECT COUNT(*) OUTPUT FROM mapping`
	var count int
	err = db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		logging.Logger.Errorf("Error querying row count: %v", err)
		logging.Logger.Info("Attempting to initialize Database")
		PostgresInit()
	}

	if len(jsonResp) == count {
		return
	}

	// Log when attempting to add items
	logging.Logger.Infof("Adding %d item(s) to Database...", len(jsonResp)-count)

	// Upsert all item data
	for _, item := range jsonResp {
		// Gather required data for insertion
		tempObj := models.Item{
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

		// Replace icon url with encoded url for dynamic image fetching
		tempObj.Icon = encodeURL(tempObj.Icon.(string))

		// Create and Execute INSERT statement
		insertQuery := `INSERT INTO public.mapping (id, members, lowalch, highalch, buylimit, value, icon, name, examine)
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
			logging.Logger.Fatalf("Failed to execute insert query: %v", err)
		}
	}

	db.Close()

	logging.Logger.Info("Successfully added item updates.")
}

func encodeURL(imageURL string) string {
	encoded := strings.ReplaceAll(imageURL, " ", "_")
	return encoded
}
