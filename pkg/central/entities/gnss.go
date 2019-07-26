package entities

import (
	"fmt"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/sirupsen/logrus"
	"github.com/stratoberry/go-gpsd"
)

type GNSSData struct {
	Location       *Location  `json:"location"`
	SatelliteCount *int64     `json:"satellite_count"`
	FixQuality     *gpsd.Mode `json:"fix_quality"`
	Time           *time.Time `json:"time"`
	GroundSpeedKPH *float64   `json:"speed_kph"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SatelliteFix struct {
	Location
	SatelliteCount int64
	FixQuality     string
}

func (gnss *GNSSData) UpdateFromGLL(gll nmea.GLL) {
	gnss.Location = &Location{
		Latitude:  gll.Latitude,
		Longitude: gll.Longitude,
	}
}

func (gnss *GNSSData) UpdateFromGGA(gga nmea.GGA) {
	// Prep
	location := Location{
		Latitude:  gga.Latitude,
		Longitude: gga.Longitude,
	}
	fixQuality := gpsd.Mode(0)
	if len(gga.FixQuality) > 0 {
		fixQuality = gpsd.Mode(gga.FixQuality[0]) // FIXME - this is probably dangerous
	}
	satelliteCount := gga.NumSatellites

	// Update
	gnss.Location = &location
	gnss.FixQuality = &fixQuality
	gnss.SatelliteCount = &satelliteCount
}

func (gnss *GNSSData) UpdateFromVTG(vtg nmea.VTG) {
	gnss.GroundSpeedKPH = &vtg.GroundSpeedKPH
}

func (gnss *GNSSData) UpdateFromZDA(zda nmea.ZDA) {
	// Prep
	// Examples: 2002-10-02T10:00:00-05:00; 2006-01-02T15:04:05Z07:00
	date := fmt.Sprintf("%04d-%02d-%02d", zda.Year, zda.Month, zda.Day)                          // 2006-01-02
	hourMinSec := fmt.Sprintf("%02d:%02d:%02d", zda.Time.Hour, zda.Time.Minute, zda.Time.Second) // 15:04:05
	localOffset := fmt.Sprintf("%02d:%02d", zda.OffsetHours, zda.OffsetMinutes)                  // 07:00

	zdaTimeValue := fmt.Sprintf("%sT%sZ%s", date, hourMinSec, localOffset)
	zdaTime, err := time.Parse(time.RFC3339, zdaTimeValue) // 2006-01-02T15:04:05Z07:00
	if err != nil {
		// No-op if err
		logrus.Errorln("Unable to parse date/time from ZDA:", err)
		return
	}

	// Update
	gnss.Time = &zdaTime
}

func (gnss *GNSSData) UpdateFromSKY(sky gpsd.SKYReport) {
	// Prep
	numSats := int64(0)
	for _, sat := range sky.Satellites {
		if sat.Used {
			numSats++
		}
	}

	// Update
	gnss.SatelliteCount = &numSats
	gnss.Time = &sky.Time
}

func (gnss *GNSSData) UpdateFromTPV(tpv gpsd.TPVReport) {
	// Prep
	location := Location{
		Latitude:  tpv.Lat,
		Longitude: tpv.Lon,
	}

	// Update
	gnss.Location = &location
	gnss.FixQuality = &tpv.Mode
	gnss.Time = &tpv.Time
}
