package main

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/parrotmac/rusted/pkg/central/services"
	"github.com/parrotmac/rusted/pkg/device/cellular"
	"github.com/parrotmac/rusted/pkg/device/gnss"
	"github.com/parrotmac/rusted/pkg/transport"
)

type LocationProviderName string

const (
	LocationProviderSerialNmea LocationProviderName = "serial-nmea"
	LocationProviderGpsd       LocationProviderName = "gpsd"
)

var ErrNoLocationProviderSpecified = errors.New("no location provider set")

type Rusted struct {
	ClientID             string
	LocationProviderName LocationProviderName
	*services.LocationService
	*services.CellularService
	*transport.MqttWrapper
}

func (r *Rusted) initContext() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "__INVALID"
	}
	r.ClientID = hostname
}

func getLocationProvider() (LocationProviderName, error) {
	providerName := os.Getenv("RUSTED_LOCATION_PROVIDER_NAME")
	if providerName == string(LocationProviderSerialNmea) {
		return LocationProviderSerialNmea, nil
	}
	if providerName == string(LocationProviderGpsd) {
		return LocationProviderGpsd, nil
	}
	return "", ErrNoLocationProviderSpecified
}

func (r *Rusted) initApp() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// FIXME
		panic(err)
	}
	// TODO: Switch to Zap
	logrus.Debugln("[APP INIT] Starting Init")

	cellInfo := &cellular.StatusProvider{
		Logger: logger,
	}
	err = cellInfo.Setup()
	if err != nil {
		logrus.Fatalf("Failed to setup cell info service: %v", err)
	}

	mqttWrapper, err := transport.ConnectMqttWrapper(getMqttConfig())
	if err != nil {
		logrus.Fatalln(err)
	}
	r.MqttWrapper = mqttWrapper

	providerName, err := getLocationProvider()
	if err != nil {
		// Setting the default to GPSD
		logrus.Warnln("no location provider explicitly set, defaulting to gpsd")
		providerName = LocationProviderGpsd
	}

	if providerName == LocationProviderSerialNmea {
		gnssReceiver, err := gnss.StartReceiver("/dev/ttyACM0", 115200)
		if err != nil {
			logrus.Fatalf("Failed to open GPS receiver: %v", err)
		}
		r.LocationService = services.CreateLocationReportingService(gnssReceiver, r.MqttWrapper)
	}
	if providerName == LocationProviderGpsd {
		r.LocationService = services.CreateLocationReportingService(&gnss.GPSDClient{}, r.MqttWrapper)
	}

	if r.LocationService != nil {
		go r.LocationService.Start()
	}

	r.CellularService = services.CreateCellularService(time.Second*10, cellInfo, r.MqttWrapper)
	r.CellularService.Start()
}

func getHttpConfig() transport.HttpConfig {
	return transport.HttpConfig{
		SeverBaseURL:   "https://api.example.com:8080",
		DefaultTimeout: time.Duration(time.Second * 5),
	}
}

func getMqttConfig() transport.MqttConfig {
	return transport.MqttConfig{
		BrokerURL: "tcp://mqtt.stag9.com:1883",
	}
}

func (r *Rusted) runLoop() {
	logrus.Warnln("[MAIN LOOP] Not yet implemented")
	for {
		time.Sleep(time.Second + 1)
	}
}

func main() {
	logrus.Info("Starting Rusted")

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	r := &Rusted{}
	r.initContext()

	r.initApp()

	r.runLoop()

	logrus.Warnln("Exiting...")
}
