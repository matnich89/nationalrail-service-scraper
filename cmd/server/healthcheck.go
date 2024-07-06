package cmd

import (
	"log"
	"net/http"
)

func (a *App) AddHealthCheckEndpoint() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Println("error writing healthcheck response")
		}
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("Health check server failed: %v", err)
		}
	}()
}
