package status

import "time"

type PeriodicUpdate struct {
	LastUpdatedAt    time.Time
	FrequencySeconds int
}

type StatusReportingConfig struct {
	ModemStatus PeriodicUpdate
	Location    PeriodicUpdate
}
