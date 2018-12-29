package server

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

func (r *Remote) ConnectMqttWrapper(mqttBrokerURL string) *MqttWrapper {
	_, wrapper := connectMqttWrapper(mqttBrokerURL)
	return wrapper
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
