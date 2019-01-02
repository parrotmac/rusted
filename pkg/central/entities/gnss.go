package entities

import "github.com/adrianmo/go-nmea"

type BasicLocation struct {
	// Filled from GLL data
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Time      Time   `json:"time"`
}

type AdvancedLocation struct {
	// Filled from GGA data
	Latitude       string  `json:"latitude"`
	Longitude      string  `json:"longitude"`
	SatelliteCount int64   `json:"sat_count"`
	Altitude       float64 `json:"altitude"`
	FixQuality     string  `json:"fix_quality"`
	Time           Time    `json:"time"`
}

type GroundSpeed struct {
	SpeedKPH float64 `json:"speed_kph"`
}

// Basically just go-nmea's location
type Time struct {
	Hour        int `json:"h"`
	Minute      int `json:"m"`
	Second      int `json:"s"`
	Millisecond int `json:"ms"`
}

func timeFromNmeaTime(time nmea.Time) Time {
	return Time{
		Hour:        time.Hour,
		Minute:      time.Minute,
		Second:      time.Second,
		Millisecond: time.Millisecond,
	}
}

func NewBasicLocationFromGLL(gll nmea.GLL) BasicLocation {
	return BasicLocation{
		Latitude:  nmea.FormatGPS(gll.Latitude),
		Longitude: nmea.FormatGPS(gll.Longitude),
		Time:      timeFromNmeaTime(gll.Time),
	}
}

func NewAdvancedLocationFromGGA(gga nmea.GGA) AdvancedLocation {
	return AdvancedLocation{
		Latitude:       nmea.FormatGPS(gga.Latitude),
		Longitude:      nmea.FormatGPS(gga.Longitude),
		SatelliteCount: gga.NumSatellites,
		Altitude:       gga.Altitude,
		Time:           timeFromNmeaTime(gga.Time),
	}
}

func NewGroundSpeedFromVTG(vgg nmea.VTG) GroundSpeed {
	return GroundSpeed{
		SpeedKPH: vgg.GroundSpeedKPH,
	}
}
