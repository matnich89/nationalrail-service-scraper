package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/matnich89/national-rail-client/nationalrail"
	"log"
	"sync"
	"time"
	"trainstats-scraper/model"
)

type Worker struct {
	ID           int
	Stations     []nationalrail.CRSCode
	ServiceIDS   map[string]bool
	ServiceChan  chan model.Train
	Client       *nationalrail.Client
	InitialDelay time.Duration
	Ticker       *time.Ticker
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
	departureBoard, err := w.Client.GetDepartures(station)
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
	if _, ok := w.ServiceIDS[service.ID]; !ok {
		if service.ScheduledTimeOfDeparture != nil {
			scheduledTime, err := parseScheduledTime(*service.ScheduledTimeOfDeparture)
			if err != nil {
				return fmt.Errorf("could not parse scheduled time: %w", err)
			}

			train := model.Train{
				ID:                 service.ID,
				ScheduledDeparture: scheduledTime,
			}

			select {
			case w.ServiceChan <- train:
				w.ServiceIDS[service.ID] = true
			case <-ctx.Done():
				return context.Canceled
			}
		}
	} else {
		log.Println("already have service", service.ID)
	}
	return nil
}

func parseScheduledTime(scheduledTime string) (time.Time, error) {
	t, err := time.Parse("15:04", scheduledTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
