package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parrotmac/rusted/pkg/central"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type Remote struct {
	ReportTypes  ReportTypes
	CommandTypes CommandTypes
}

type Report struct {
	// Used as suffix to mqttClient pub topic
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
		TypeID: "/loc/detail",
	}

	r.ReportTypes.CellCarrier = Report{
		TypeID: "/cell/carrier",
	}

	r.ReportTypes.CellQuality = Report{
		TypeID: "/cell/quality",
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

/*
This might turn into a goroutine that connects and manages reconnections independently
*/
// This is a goroutine that's called
func (w *MqttWrapper) maintainMqttConnection() {
	retries := -1 // First connection is free

	for {
		if !(*w.mqttClient).IsConnected() {
			retries++
		}
	}
}

func ConnectMqttWrapper(cfg MqttConfig) (*MqttWrapper, error) {
	err, wrapper := connectMqttWrapper(cfg.BrokerURL)
	if err != nil {
		logrus.Warnf("Unable to connect to MQTT broker: %v", err)
		return nil, err
	}
	return wrapper, nil
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

func (w *MqttWrapper) publishToTopic(topic string, data interface{}) error {
	dataPayload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if w.mqttClient != nil && (*w.mqttClient).IsConnected() {
		(*w.mqttClient).Publish(topic, 1, false, dataPayload)
	} else {
		return errors.New("mqttClient is nil or is not connected")
	}
	return nil
}

func (w *MqttWrapper) ReportBasicLocation(ctx *central.Context, location entities.BasicLocation) error {
	topic := fmt.Sprintf("evt/%s/loc/basic", ctx.ClientIdentifier)
	err := w.publishToTopic(topic, location)
	if err != nil {
		logrus.Warnf("Unable to publish: %v", err)
		return err
	}
	return nil
}

func (w *MqttWrapper) ReportDetailedLocation(ctx *central.Context, location entities.AdvancedLocation) error {
	topic := fmt.Sprintf("evt/%s/loc/detail", ctx.ClientIdentifier)
	err := w.publishToTopic(topic, location)
	if err != nil {
		logrus.Warnf("Unable to publish: %v", err)
		return err
	}
	return nil
}

func (w *MqttWrapper) ReportCellQuality(ctx *central.Context, quality entities.CellQuality) error {
	topic := fmt.Sprintf("evt/%s/cell/quality", ctx.ClientIdentifier)
	err := w.publishToTopic(topic, quality)
	if err != nil {
		logrus.Warnf("Unable to publish: %v", err)
		return err
	}
	return nil
}
