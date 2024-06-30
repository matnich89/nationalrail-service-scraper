package internal

import (
	"github.com/matnich89/national-rail-client/nationalrail"
	"log"
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

func (w *Worker) Work(done <-chan struct{}) {
	time.Sleep(w.InitialDelay)

	log.Printf("worker %d starting...", w.ID)

	defer w.Ticker.Stop()

	w.checkStations()

	for {
		select {
		case <-w.Ticker.C:
			w.checkStations()
		case <-done:
			return
		}
	}
}

func (w *Worker) checkStations() {
	for _, station := range w.Stations {
		departureBoard, err := w.Client.GetDepartures(station)
		if err != nil {
			log.Printf("worker %d could not get departure board for %s, error is %v", w.ID, station, err)
			continue
		}
		services := departureBoard.Services
		if services == nil || len(services) == 0 {
			log.Printf("no services currently scheduled at station %s", station)
		}
		for _, service := range services {
			if _, ok := w.ServiceIDS[service.ID]; !ok {
				if service.ScheduledTimeOfDeparture != nil {
					scheduledTime, err := parseScheduledTime(*service.ScheduledTimeOfDeparture)
					if err != nil {
						log.Printf("could not parse scheduled time for station %s, error is %v", station, err)
						continue
					}
					train := model.Train{
						ID:                 service.ID,
						ScheduledDeparture: scheduledTime,
					}
					w.ServiceChan <- train
					w.ServiceIDS[service.ID] = true
				}
			} else {
				log.Println("already have service", station)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func parseScheduledTime(scheduledTime string) (time.Time, error) {
	t, err := time.Parse("15:04", scheduledTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
