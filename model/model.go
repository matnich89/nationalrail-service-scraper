package model

import "time"

type Train struct {
	ID                 string
	ScheduledDeparture time.Time
}
