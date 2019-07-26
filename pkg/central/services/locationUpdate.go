package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type LocationService struct {
	UpdateFrequency time.Duration
	entities.LocationProvider
	entities.LocationUpdateTransport
}

func CreateLocationReportingService(
	locProvider entities.LocationProvider,
	transport entities.LocationUpdateTransport,
) *LocationService {
	return &LocationService{
		UpdateFrequency:         time.Duration(time.Second * 10),
		LocationProvider:        locProvider,
		LocationUpdateTransport: transport,
	}
}

func (srv *LocationService) handleLocationUpdate(data entities.GNSSData) {
	err := srv.LocationUpdateTransport.SendLocationReport(context.TODO(), data)
	if err != nil {
		logrus.Warnln("failure reporting location update:", err)
	}
}

func (srv *LocationService) Start() {
	// Start reporting location
	srv.LocationProvider.SetOnLocationUpdatedHandler(srv.handleLocationUpdate)
	err := srv.LocationProvider.Start()
	if err != nil {
		logrus.Fatal("failed to start location provider:", err)
	}
}
