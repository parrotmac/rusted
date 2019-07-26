package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type MqttConfig struct {
	BrokerURL string
}

type MqttWrapper struct {
	MqttConfig       *MqttConfig
	mqttClient       *mqtt.Client
	onConnect        func()
	ATCommandHandler func(atCommand string) string
	ClientID         string
}

func connectMqttWrapper(brokerURL string) (error, *MqttWrapper) {
	wrapper := MqttWrapper{
		ClientID: "7dfa82ff-1989-44a7-a9fd-befddfb93ad9",
	}
	opts := &mqtt.ClientOptions{}

	opts.AddBroker(brokerURL)
	opts.CleanSession = true
	opts.ClientID = wrapper.ClientID
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
		logrus.Infoln(string(message.Payload()))
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
			logrus.Warnln("Got empty payload :(")
		}
	})
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

func (w *MqttWrapper) SendLocationReport(ctx context.Context, location entities.GNSSData) error {
	topic := fmt.Sprintf("evt/%s/loc", w.ClientID)
	err := w.publishToTopic(topic, location)
	if err != nil {
		logrus.Warnf("Unable to publish: %v", err)
		return err
	}
	return nil
}

func (w *MqttWrapper) SendCellInfoReport(ctx context.Context, report *entities.ModemReport) error {
	topic := fmt.Sprintf("evt/%s/cell", w.ClientID)
	err := w.publishToTopic(topic, report)
	if err != nil {
		logrus.Warnf("Unable to publish: %v", err)
		return err
	}
	return nil
}

func (w *MqttWrapper) ReportCellQuality(ctx context.Context, quality entities.ModemReport) error {
	topic := fmt.Sprintf("evt/%s/cell/quality", w.ClientID)
	err := w.publishToTopic(topic, quality)
	if err != nil {
		logrus.Warnf("Unable to publish: %v", err)
		return err
	}
	return nil
}
