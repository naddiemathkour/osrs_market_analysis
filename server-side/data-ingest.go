package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello World")

	//Handle http request
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/latest", nil)
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
	}

	//Decode response into pretty print json
	var jsonResp map[string]interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the JSON with indentation
	// prettyJson, err := json.MarshalIndent(jsonResp, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//Iterate over the data
	for id, item := range jsonResp {
		fmt.Println(id, item)
		fmt.Println("cringe")
	}

	// fmt.Println(string(prettyJson))
}
