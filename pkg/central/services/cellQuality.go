package services

import (
	"time"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type CellularService struct {
	UpdateFrequency time.Duration
}

func CreateCellularService(locProvider entities.LocationProvider, transport entities.LocationUpdateTransport) CellularService {
	return CellularService{
		UpdateFrequency: time.Duration(time.Second * 30),
	}
}
