package serv

import (
	"net/http"

	db "github.com/naddiemathkour/osrs_market_analysis/db"
	"github.com/naddiemathkour/osrs_market_analysis/logging"
)

func StartServer() *http.Server {
	// Create routing mux
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", corsMiddleware(db.DataFetchHandler))

	// Initialize Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Create go routine to start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatalf("Http server failed to start: %v", err)
		}
	}()

	logging.Logger.Info("HTTP server started on :8080")
	return srv
}
