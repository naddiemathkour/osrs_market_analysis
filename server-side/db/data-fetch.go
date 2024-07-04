package db

import (
	"fmt"
	"net/http"
)

func DataFetchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Printing from DataFetchHandler")
}
