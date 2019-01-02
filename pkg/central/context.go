package central

import (
	"github.com/eclipse/paho.mqtt.golang"
	"net/http"
)

type Context struct {
	HttpConfig struct{}    // TODO
	HttpClient http.Client // TODO: Use wrapper

	MqttConfig struct{}    // TODO:
	MqttClient mqtt.Client // TODO Use wrapper

	SmsConfig struct{}    // TODO
	SmsClient interface{} // TODO

	// For now all reporting is tied to the same frequency
	StatusReportingFrequency int
}
