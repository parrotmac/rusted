package central

import "github.com/parrotmac/rusted/pkg/central/entities"

type Capabilities interface {
	/* Location & GPS Services */
	GetBasicLocation(ctx Context) entities.BasicLocation
	GetAdvancedLocation(ctx Context) entities.AdvancedLocation
	GetGPSReportedSpeed(ctx Context) entities.GroundSpeed

	// It might make sense to combine these into a single message
	GetCellSignalQuality(ctx Context) entities.CellQuality
	GetCellCarrierName(ctx Context) entities.CellCarrier

	// Sends to default server
	SendSMSMessage(ctx Context, payload []byte)

	// Update

	// May not be used right away but demonstrates the pattern
	RegisterGPSDataReceived(ctx Context, gpsLocationChanged func() entities.AdvancedLocation)
}
