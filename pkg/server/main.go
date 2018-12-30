package server

import (
	"encoding/json"
	"github.com/parrotmac/rusted/pkg/gnss"
	"github.com/sirupsen/logrus"
)

type Remote struct {
	httpBaseURL string

	mqttConnectionValid bool
	mqttBrokerURL       string
	mqttWrapper         *MqttWrapper
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

func (r *Remote) PublishBasicLocationUpdate(location gnss.BasicLocation) {
	locationData, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to marshal JSON message: %v", err)
	}
	if r.mqttWrapper != nil && r.mqttWrapper.mqttClient != nil && (*r.mqttWrapper.mqttClient).IsConnected() {
		(*r.mqttWrapper.mqttClient).Publish("event/yogurt/location/basic", 1, false, locationData)
	}
}

func (r *Remote) PublishAdvancedLocationUpdate(location gnss.AdvancedLocation) {
	locationData, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to marshal JSON message: %v", err)
	}
	if r.mqttWrapper != nil && r.mqttWrapper.mqttClient != nil && (*r.mqttWrapper.mqttClient).IsConnected() {
		(*r.mqttWrapper.mqttClient).Publish("event/yogurt/location/advanced", 1, false, locationData)
	}
}
