package db

import (
	"encoding/json"
	"net/http"

	"github.com/naddiemathkour/osrs_market_analysis/logging"
)

type DataResponse struct {
	ItemListings []map[string]interface{} `json:"items"`
}

func DataFetchHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch DB connection object
	db := Connect(r.Method)

	// Verify the connection
	if err := db.Ping(); err != nil {
		logging.Logger.Fatalf("Failed to ping database: %v", err)
	}

	queryString := `SELECT * FROM Listings
					ORDER BY (margin) desc;`

	rows, err := db.Queryx(queryString)
	if err != nil {
		logging.Logger.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	// Process query result
	var items []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			logging.Logger.Fatalf("Failed to scan row: %v", err)
		}
		items = append(items, row)
	}

	// Check for errors after iterating over rows
	if err = rows.Err(); err != nil {
		logging.Logger.Fatalf("Error iterating over rows: %v", err)
	}

	response := DataResponse{
		ItemListings: items,
	}

	db.Close()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
