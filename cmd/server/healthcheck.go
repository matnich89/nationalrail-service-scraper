package cmd

import (
	"log"
	"net/http"
)

func (a *App) AddHealthCheckEndpoint() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("Health check server failed: %v", err)
		}
	}()
}
