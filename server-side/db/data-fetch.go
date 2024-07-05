package db

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DataResponse struct {
	Message string `json:"message"`
}

func DataFetchHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch DB connection object
	db := Connect(r.Method)

	fmt.Println(db)

	response := DataResponse{
		Message: "Hello",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
