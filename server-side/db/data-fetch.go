package db

import (
	"fmt"
	"net/http"
)

func DataFetchHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch DB connection object
	db := Connect(r.Method)

	fmt.Println(db)
}
