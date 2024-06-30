package cmd

import (
	"context"
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

type App struct {
	client     *nr.Client
	numWorkers int
	stations   []nr.CRSCode
	workers    []*internal.Worker
	wg         sync.WaitGroup
	trainChan  chan model.Train
}

func NewApp(client *nr.Client, numWorkers int, stations []nr.CRSCode) *App {
	return &App{
		client:     client,
		numWorkers: numWorkers,
		stations:   stations,
		wg:         sync.WaitGroup{},
		trainChan:  make(chan model.Train),
	}
}

func (a *App) SetupWorkers() {

	log.Println("setting up workers")

	maxDelay := 3 * time.Second
	stationsPerWorker := len(a.stations) / a.numWorkers

	for i := 0; i < a.numWorkers; i++ {
		start := i * stationsPerWorker
		end := start + stationsPerWorker
		if i == a.numWorkers-1 {
			end = len(a.stations)
		}

		worker := &internal.Worker{
			ID:           i,
			Stations:     a.stations[start:end],
			ServiceIDS:   make(map[string]bool),
			ServiceChan:  a.trainChan,
			Client:       a.client,
			InitialDelay: time.Duration(rand.Int63n(int64(maxDelay))), // to stop all the workers hitting national rail at once, stagger the initial invocation
			Ticker:       time.NewTicker(5 * time.Minute),             // departure boards are fairly slow to change, so can just check every 5 mins, to pick up new services
		}

		a.workers = append(a.workers, worker)
	}

	log.Printf("created %d workers", len(a.workers))

}

func (a *App) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("Received signal %v. Starting shutdown...\n", sig)
		cancel()
	}()

	// start workers
	for _, worker := range a.workers {
		a.wg.Add(1)
		go worker.Work(ctx, &a.wg)
	}

	go a.listen()

	a.wg.Wait()
	close(a.trainChan)
}

func (a *App) listen() {
	for train := range a.trainChan {
		log.Printf("received train %s", train.ID)
	}
}
