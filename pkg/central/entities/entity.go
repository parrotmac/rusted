package entities

import (
	"context"
)

// MQTT and SMS should both implement this interface
type LocationUpdateTransport interface {
	SendLocationReport(ctx context.Context, locationData GNSSData) error
}

type CellInfoTransport interface {
	SendCellInfoReport(ctx context.Context, report *ModemReport) error
}

// Serial GPS receivers, or a gpsd client
type LocationProvider interface {
	SetOnLocationUpdatedHandler(func(locationData GNSSData))
	Start() error
}
