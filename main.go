package main

import (
	"fmt"
	nr "github.com/matnich89/national-rail-client/nationalrail"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"trainstats-scraper/internal"
	"trainstats-scraper/model"
)

var (
	stationCheckedCounter *internal.StationCheckCounter
)

func main() {
	nrClient, err := nr.NewClient(
		nr.AccessTokenOpt(os.Getenv("NATIONAL_RAIL_API_KEY")),
	)

	if err != nil {
		log.Fatalf("could not create national rail client: %v", err)
	}

	numWorkers := 20
	stations := internal.GetStations()
	fmt.Printf("there are %d stations to check", len(stations))
	stationsPerWorker := len(stations) / numWorkers

	done := make(chan struct{})
	var wg sync.WaitGroup

	maxDelay := time.Minute * 1

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serviceChan := make(chan model.Train, 1000)

	for i := 0; i < numWorkers; i++ {
		start := i * stationsPerWorker
		end := start + stationsPerWorker
		if i == numWorkers-1 {
			end = len(stations)
		}

		worker := &internal.Worker{
			ID:           i,
			Stations:     stations[start:end],
			ServiceIDS:   make(map[string]bool),
			ServiceChan:  serviceChan,
			Client:       nrClient,
			InitialDelay: time.Duration(rand.Int63n(int64(maxDelay))), // to stop all the workers hitting national rail at once, stagger the initial invocation
			Ticker:       time.NewTicker(5 * time.Minute),             // departure boards are fairly slow to change, so can just check every 5 mins, to pick up new services
		}

		wg.Add(1)

		go func(w *internal.Worker) {
			defer wg.Done()
			w.Work(done)
		}(worker)
	}

	go func() {
		for service := range serviceChan {
			// TODO push to queue
			fmt.Printf("Received service: %v\n ", service)
		}
	}()

	wg.Wait()

	close(serviceChan)

	fmt.Println("Shutting down...")

}
