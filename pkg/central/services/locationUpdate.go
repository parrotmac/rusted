package services

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central"
	"github.com/parrotmac/rusted/pkg/central/entities"
)

type LocationService struct {
	*central.Context
	UpdateFrequency time.Duration
	entities.LocationProvider
	entities.LocationUpdateTransport
}

func CreateLocationReportingService(
	ctx *central.Context,
	locProvider entities.LocationProvider,
	transport entities.LocationUpdateTransport,
) LocationService {
	return LocationService{
		Context:                 ctx,
		UpdateFrequency:         time.Duration(time.Second * 10),
		LocationProvider:        locProvider,
		LocationUpdateTransport: transport,
	}
}

func (srv *LocationService) sendUpdate() error {
	basicLocation, err := srv.GetBasicLocation(srv.Context)
	if err != nil {
		logrus.Debugf("Failed to get basic location: %v", err)
		return err
	}

	err = srv.ReportBasicLocation(srv.Context, basicLocation)
	if err != nil {
		logrus.Debugf("Failed to report location update: %v", err)
		return err
	}

	return nil
}

func (srv *LocationService) RunService() {
	// Start reporting location
	err := srv.LocationProvider.Start(srv.Context)
	if err != nil {
		// TODO: Bubble to smart retry logic
		logrus.Warnf("Failed to start location service: %v", err)
		logrus.Fatalln("Exiting as this should be handled differently")
	}
	for {
		err := srv.sendUpdate()
		if err != nil {
			// TODO: Modify flow based on error type and possibly attempt to resolve or notify of failure
		}

		// TODO: Store update time and wait for changes or time to pass
		time.Sleep(srv.UpdateFrequency)
	}
}
