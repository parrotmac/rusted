package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/central/entities"
	"github.com/parrotmac/rusted/pkg/device/cellular"
)

type SMSConfig struct {
	CellModem interface{} // TODO: Define and build this
}

type SmsClient struct {
	ClientID string
	cellular.SmsService
}

func (s *SmsClient) ReportBasicLocation(ctx context.Context, location entities.GNSSData) error {
	// TODO: Do we want to include MQTT topics in SMS messages?
	mqttTopic := fmt.Sprintf("evt/%s/loc/basic", s.ClientID)
	messageBody, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	message := cellular.CreateOutgoingSMS("9999999999", []byte(fmt.Sprintf("%s\n%s", mqttTopic, messageBody)))
	err = s.SmsService.SendMessage(message)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	return nil
}

func (s *SmsClient) ReportDetailedLocation(ctx context.Context, location entities.GNSSData) error {
	// TODO: Do we want to include MQTT topics in SMS messages?
	mqttTopic := fmt.Sprintf("evt/%s/loc/detail", s.ClientID)
	messageBody, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	message := cellular.CreateOutgoingSMS("9999999999", []byte(fmt.Sprintf("%s\n%s", mqttTopic, messageBody)))
	err = s.SmsService.SendMessage(message)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	return nil
}

func (s *SmsClient) ReportCellQuality(ctx context.Context, quality entities.ModemReport) error {
	return errors.New("not implemented")
}
