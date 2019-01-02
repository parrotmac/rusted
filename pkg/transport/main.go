package transport

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type Remote struct {
	httpBaseURL string

	mqttConnectionValid bool
	mqttBrokerURL       string
	mqttWrapper         *MqttWrapper

	deviceIdentifier string

	ReportTypes  ReportTypes
	CommandTypes CommandTypes
}

type Report struct {
	// Used as suffix to mqtt pub topic
	// for 'Basic Location' the `TypeID` would be '/loc/basic'
	TypeID string
}

type ReportTypes struct {
	BasicLocation    Report
	DetailedLocation Report
	CellQuality      Report
	CellCarrier      Report
}

// Hopefully abstract enough to support MQTT or SMS
type CommandHandler func(commandPath string, commandBody []byte)

type Command struct {
	TypeID         string
	CommandHandler CommandHandler
}

type CommandTypes struct {
	SetReportingFrequency Command
	DoorLockActuation     Command
	TrunkPop              Command
	RemoteStart           Command
	FastHonk              Command
}

func (r *Remote) SetupReporting() {
	r.ReportTypes.BasicLocation = Report{
		TypeID: "/loc/basic",
	}

	r.ReportTypes.DetailedLocation = Report{
		TypeID: "loc/detail",
	}

	r.ReportTypes.CellCarrier = Report{
		TypeID: "cell/carrier",
	}

	r.ReportTypes.CellQuality = Report{
		TypeID: "cell/quality",
	}
}

func defaultCommandHandler(commandPath string, commandBody []byte) {
	logrus.Warnf("Unhandled command received: %v/%v", commandPath, commandBody)
}

func (r *Remote) SetupCommandReceivers() {
	r.CommandTypes.SetReportingFrequency = Command{
		TypeID:         "/report/set-freq/+",
		CommandHandler: defaultCommandHandler,
	}

	r.CommandTypes.FastHonk = Command{
		TypeID:         "/trick/fasthonk",
		CommandHandler: defaultCommandHandler,
	}

	r.CommandTypes.RemoteStart = Command{
		TypeID:         "/engine/remote",
		CommandHandler: defaultCommandHandler,
	}
}

func NewRemote() *Remote {
	return &Remote{}
}

/*
This might turn into a goroutine that connects and manages reconnections independently
*/
// This is a goroutine that's called
func (r *Remote) maintainMqttConnection() {
	retries := -1 // First connection is free

	for {
		if !(*r.mqttWrapper.mqttClient).IsConnected() {

			retries++
		}
	}
}

func (r *Remote) ConnectMqttWrapper(mqttBrokerURL string) error {
	err, wrapper := connectMqttWrapper(mqttBrokerURL)
	if err != nil {
		logrus.Warnf("Unable to connect to MQTT broker: %v", err)
		return err
	}
	r.mqttWrapper = wrapper
	return nil
}

func (r *Remote) GetMqttConnectionIsValid() bool {
	// Other components may wish to check the health of the connection
	return r.mqttConnectionValid
}

/*****************************************************
*
*   These provide a way to publish update messages, but don't expose details of transport mechanism
*   It probably makes sense to move these either to a separate package (or even just a separate file), and/or to define
*   an interface that the MQTT, HTTP, SMS, etc. remote connections could implement. However, it might make sense for
*   those transport mechanisms to simply provide generic read/write/notify/etc. handlers that an telemetry/command
*   system might utilize
*
 *****************************************************/
func (r *Remote) PublishCarrierStatus(carrier string) error {
	return nil
}

func (r *Remote) PublishSignalStrengthStatus(signalDbm string) error {
	return nil
}

func (r *Remote) PublishBasicLocationUpdate(location entities.BasicLocation) {
	locationData, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to marshal JSON message: %v", err)
	}
	if r.mqttWrapper != nil && r.mqttWrapper.mqttClient != nil && (*r.mqttWrapper.mqttClient).IsConnected() {
		(*r.mqttWrapper.mqttClient).Publish("event/yogurt/location/basic", 1, false, locationData)
	}
}

func (r *Remote) PublishAdvancedLocationUpdate(location entities.AdvancedLocation) {
	locationData, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to marshal JSON message: %v", err)
	}
	if r.mqttWrapper != nil && r.mqttWrapper.mqttClient != nil && (*r.mqttWrapper.mqttClient).IsConnected() {
		(*r.mqttWrapper.mqttClient).Publish("event/yogurt/location/advanced", 1, false, locationData)
	}
}
