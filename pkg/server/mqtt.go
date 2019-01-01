package server

import (
	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MqttWrapper struct {
	mqttClient       *mqtt.Client
	onConnect        func()
	ATCommandHandler func(atCommand string) string
}

func connectMqttWrapper(brokerURL string) (error, *MqttWrapper) {
	wrapper := MqttWrapper{}
	opts := &mqtt.ClientOptions{}

	opts.AddBroker(brokerURL)
	opts.CleanSession = true
	opts.ClientID = "7dfa82ff-1989-44a7-a9fd-befddfb93ad9"
	opts.OnConnect = wrapper.attachSubscriptions

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error(), nil
	}

	wrapper.mqttClient = &client

	return nil, &wrapper
}

func (w *MqttWrapper) attachSubscriptions(client mqtt.Client) {
	client.Subscribe("cmd/yogurt/status", 1, func(client mqtt.Client, message mqtt.Message) {
		log.Infoln(string(message.Payload()))
	})
	//if token.Wait() && token.Error() != nil {
	//
	//}
	client.Subscribe("cmd/yogurt/at-cmd", 2, func(client mqtt.Client, message mqtt.Message) {
		payload := string(message.Payload())
		if payload != "" {
			resp := w.ATCommandHandler(payload)
			if resp != "" {
				cli := *w.mqttClient
				cli.Publish("reply/yogurt/at-cmd", 2, true, resp)
			}
		} else {
			log.Warnln("Got empty payload :(")
		}
	})
}
