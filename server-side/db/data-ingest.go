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
	Limit    interface{} `json:"limit"`
	Icon     interface{} `json:"icon"`
	Examine  interface{} `json:"examine"`
}

func MapItems(db *sqlx.DB) {
	//handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/mapping", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("User-Agent", "Runescape Market Data Analysis")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Decode response data to JSON
	var jsonResp []map[string]interface{}

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range jsonResp {
		tempObj := MapObject{
			ID:       item["id"],
			Name:     item["name"],
			Members:  item["members"],
			Highalch: item["highalch"],
			Lowalch:  item["lowalch"],
			Value:    item["value"],
			Limit:    item["limit"],
			Icon:     item["icon"],
			Examine:  item["examine"],
		}
		fmt.Println(tempObj)
		fmt.Println(db.Ping())
	}
}
