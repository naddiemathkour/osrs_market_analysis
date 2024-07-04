package db

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
)

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
