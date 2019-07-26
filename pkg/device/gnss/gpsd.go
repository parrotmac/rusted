package gnss

import (
	"github.com/stratoberry/go-gpsd"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

const defaultGPSDLocation = "localhost:2947"

type GPSDClient struct {
	DaemonAddress string

	gpsdSession *gpsd.Session

	dataCache entities.GNSSData

	onLocationUpdated func(data entities.GNSSData)
}

func (client *GPSDClient) getDaemonAddress() string {
	if client.DaemonAddress == "" {
		return defaultGPSDLocation
	}
	return client.DaemonAddress
}

func (client *GPSDClient) SetOnLocationUpdatedHandler(callback func(data entities.GNSSData)) {
	client.onLocationUpdated = callback
}

func (client *GPSDClient) notifyLocationUpdated() {
	if client.onLocationUpdated != nil {
		client.onLocationUpdated(client.dataCache)
	}
}

func (client *GPSDClient) Start() error {
	gps, err := gpsd.Dial(client.getDaemonAddress())
	if err != nil {
		return err
	}
	client.gpsdSession = gps

	gps.AddFilter("TPV", func(i interface{}) {
		report := i.(*gpsd.TPVReport)
		client.dataCache.UpdateFromTPV(*report)
		client.notifyLocationUpdated()
	})

	gps.AddFilter("SKY", func(i interface{}) {
		report := i.(*gpsd.SKYReport)
		client.dataCache.UpdateFromSKY(*report)
		client.notifyLocationUpdated()
	})

	done := client.gpsdSession.Watch()
	<-done
	return nil
}
