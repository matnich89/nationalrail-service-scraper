package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/matnich89/national-rail-client/nationalrail"
	"log"
	"sync"
	"time"
	"trainstats-scraper/model"
	"trainstats-scraper/redis"
)

type Worker struct {
	ID           int
	Stations     []nationalrail.CRSCode
	ServiceChan  chan model.DepartingTrainId
	NRClient     *nationalrail.Client
	InitialDelay time.Duration
	Ticker       *time.Ticker
	RedisClient  redis.IRedisClient
}

func (w *Worker) Work(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(w.InitialDelay)

	log.Printf("worker %d starting...", w.ID)

	defer w.Ticker.Stop()

	w.checkStations(ctx)

	for {
		select {
		case <-w.Ticker.C:
			w.checkStations(ctx)
		case <-ctx.Done():
			log.Printf("worker %d stopping...", w.ID)
			return
		}
	}
}

func (w *Worker) checkStations(ctx context.Context) {
	for _, station := range w.Stations {
		stationCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := w.checkStation(stationCtx, station); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("worker %d stopping during station check...", w.ID)
				return
			}
			log.Printf("worker %d error checking station %s: %v", w.ID, station, err)
		}

		if ctx.Err() != nil {
			log.Printf("worker %d stopping during station check...", w.ID)
			return
		}
	}
}

func (w *Worker) checkStation(ctx context.Context, station nationalrail.CRSCode) error {
	departureBoard, err := w.NRClient.GetDepartures(station)
	if err != nil {
		return fmt.Errorf("could not get departure board for %s: %w", station, err)
	}

	services := departureBoard.Services
	if services == nil || len(services) == 0 {
		log.Printf("no services currently scheduled at station %s", station)
		return nil
	}

	for _, service := range services {
		if err := w.processService(ctx, service); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			log.Printf("error processing service at station %s: %v", station, err)
		}
	}

	return nil
}

func (w *Worker) processService(ctx context.Context, service *nationalrail.Service) error {

	if service.ScheduledTimeOfDeparture != nil {

		train := model.DepartingTrainId{
			ID: service.ID,
		}

		select {
		case w.ServiceChan <- train:
			// The train ID will be added to Redis in the listen function
		case <-ctx.Done():
			return context.Canceled
		}
	} else {
		log.Println("service does have scheduled time of departure", service.ID)
	}
	return nil
}
