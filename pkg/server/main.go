package server

type Remote struct {
	httpBaseURL string

	mqttBrokerURL string
	mqttWrapper *MqttWrapper
}

/*
This might turn into a goroutine that connects and manages reconnections independently
 */
func (r *Remote) ConnectMqttWrapper(mqttBrokerURL string) *MqttWrapper {
	_, wrapper := connectMqttWrapper(mqttBrokerURL)
	return wrapper
}
