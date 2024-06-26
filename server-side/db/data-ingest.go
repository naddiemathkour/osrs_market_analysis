package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
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

func MapItems(db *sqlx.DB) {
	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/mapping", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Decode response data to JSON
	var jsonResp []map[string]interface{}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}

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

		// Execute INSERT statement using sqlx.Named
		insertQuery := `INSERT INTO mapping (id, members, lowalch, highalch, buylimit, value, icon, name, examine)
                    	VALUES (:id, :members, :lowalch, :highalch, :buylimit, :value, :icon, :name, :examine)`

		fmt.Println("Inserting: ", tempObj)

		_, err := db.NamedExec(insertQuery, &tempObj)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func MapPrices(db *sqlx.DB) {
	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/5m", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Decode response data to JSON
	var jsonResp map[string]interface{}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jsonResp)

	// for _, item := range jsonResp {
	// 	// Gather required data for insertion
	// 	tempObj := MapObject{
	// 		ID:       item["id"],
	// 		Name:     item["name"],
	// 		Members:  item["members"],
	// 		Highalch: item["highalch"],
	// 		Lowalch:  item["lowalch"],
	// 		Value:    item["value"],
	// 		Buylimit: item["limit"],
	// 		Icon:     item["icon"],
	// 		Examine:  item["examine"],
	// 	}

	// 	// Execute INSERT statement using sqlx.Named
	// 	insertQuery := `INSERT INTO mapping (id, members, lowalch, highalch, buylimit, value, icon, name, examine)
	//                 	VALUES (:id, :members, :lowalch, :highalch, :buylimit, :value, :icon, :name, :examine)`

	// 	fmt.Println("Inserting: ", tempObj)

	// 	_, err := db.NamedExec(insertQuery, &tempObj)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// }
}
