package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central"
	"github.com/parrotmac/rusted/pkg/central/services"
	"github.com/parrotmac/rusted/pkg/device/gnss"
	"github.com/parrotmac/rusted/pkg/transport"
)

type Rusted struct {
	*central.Context
	*services.LocationService
	*services.CellularService
	*transport.MqttWrapper
}

func (r *Rusted) initContext() {
	r.Context = &central.Context{
		ClientIdentifier: "yogurt",
	}
}

func (r *Rusted) initApp() {
	logrus.Debugln("[APP INIT] Starting Init")

	mqttWrapper, err := transport.ConnectMqttWrapper(getMqttConfig())
	if err != nil {
		logrus.Fatalln(err)
	}
	r.MqttWrapper = mqttWrapper

	gnssReceiver, err := gnss.StartReceiver("/dev/ttyACM0", 115200)
	if err != nil {
		logrus.Warnf("Failed to open GPS receiver: %v", err)
		logrus.Warnln("No retry logic is implemented, program will now terminate")
		os.Exit(-1)
	}
	locService := services.CreateLocationReportingService(r.Context, gnssReceiver, r.MqttWrapper)
	r.LocationService = &locService

	locService.RunService()
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
