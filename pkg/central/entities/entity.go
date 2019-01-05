package entities

import "github.com/parrotmac/rusted/pkg/central"

// MQTT and SMS should both implement this interface
type LocationUpdateTransport interface {
	ReportBasicLocation(ctx *central.Context, loc BasicLocation) error
	ReportDetailedLocation(ctx *central.Context, loc AdvancedLocation) error
}

// Serial GPS receivers
type LocationProvider interface {
	GetBasicLocation(ctx *central.Context) (BasicLocation, error)
	GetDetailedLocation(ctx *central.Context) (AdvancedLocation, error)
	Start(ctx *central.Context) error
}
