package status

import "time"

type PeriodicUpdate struct {
	LastUpdatedAt    time.Time
	FrequencySeconds int
}

type Update struct {
	ModemStatus PeriodicUpdate
	Location    PeriodicUpdate
}
