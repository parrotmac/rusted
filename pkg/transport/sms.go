package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parrotmac/rusted/pkg/central"
	"github.com/parrotmac/rusted/pkg/central/entities"
	"github.com/parrotmac/rusted/pkg/device/modem"
	"github.com/sirupsen/logrus"
)

type SMSConfig struct {
	CellModem interface{} // TODO: Define and build this
}

type SmsClient struct {
	*central.Context
	modem.SmsService
}

func (s *SmsClient) ReportBasicLocation(ctx *central.Context, location entities.BasicLocation) error {
	// TODO: Do we want to include MQTT topics in SMS messages?
	mqttTopic := fmt.Sprintf("evt/%s/loc/basic", ctx.ClientIdentifier)
	messageBody, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	message := modem.CreateOutgoingSMS("9999999999", []byte(fmt.Sprintf("%s\n%s", mqttTopic, messageBody)))
	err = s.SmsService.SendMessage(s.Context, message)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	return nil
}

func (s *SmsClient) ReportDetailedLocation(ctx *central.Context, location entities.AdvancedLocation) error {
	// TODO: Do we want to include MQTT topics in SMS messages?
	mqttTopic := fmt.Sprintf("evt/%s/loc/detail", ctx.ClientIdentifier)
	messageBody, err := json.Marshal(location)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	message := modem.CreateOutgoingSMS("9999999999", []byte(fmt.Sprintf("%s\n%s", mqttTopic, messageBody)))
	err = s.SmsService.SendMessage(s.Context, message)
	if err != nil {
		logrus.Warnf("Unable to send SMS: %v", err)
		return err
	}
	return nil
}

func (s *SmsClient) ReportCellQuality(ctx *central.Context, quality entities.CellQuality) error {
	return errors.New("Not implemented")
}
