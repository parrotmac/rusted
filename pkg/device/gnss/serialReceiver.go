package gnss

import (
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"

	"github.com/parrotmac/rusted/pkg/central/entities"
)

type serialReceiver struct {
	deviceAddress string
	baudRate      int

	serialPort *serial.Port

	// TODO: Break these out
	readValueTimeout time.Duration

	gnssData entities.GNSSData

	onLocationUpdated func(data entities.GNSSData)
}

func (sr *serialReceiver) SetOnLocationUpdatedHandler(callback func(data entities.GNSSData)) {
	sr.onLocationUpdated = callback
}

func (sr *serialReceiver) notifyLocationUpdated() {
	if sr.onLocationUpdated != nil {
		sr.onLocationUpdated(sr.gnssData)
	}
}

func (sr *serialReceiver) Start() error {
	// TODO: Rethink this
	sr.processMessagesContinually()
	return nil
}

func (sr *serialReceiver) openPort() (*serial.Port, error) {

	mode := &serial.Mode{
		BaudRate: sr.baudRate,
	}
	port, err := serial.Open(sr.deviceAddress, mode)
	if err != nil {
		return nil, err
	}
	return &port, nil
}

func StartReceiver(deviceAddress string, baudRate int) (*serialReceiver, error) {
	sr := &serialReceiver{
		deviceAddress:    deviceAddress,
		baudRate:         baudRate,
		readValueTimeout: time.Duration(time.Second * 2),
	}

	port, err := sr.openPort()
	if err != nil {
		return nil, err
	}
	sr.serialPort = port

	return sr, nil
}

func (sr *serialReceiver) processMessagesContinually() {
	// TODO: This might have a main loop to pull messages from the buffer and update current information
	// This *might* use a channel or other mechanism to determine when message processing is required
	for {
		sentences, err := sr.getNextLines()
		if err != nil {
			logrus.Warnf("Read encountered error: %v", err)
			continue
		}
		for _, sentence := range sentences {
			sr.parseAndUpdate(sentence)
		}
	}
}
func (sr *serialReceiver) Run() {
	// TODO: Possibly manage connection here
	sr.processMessagesContinually()
}

func (sr *serialReceiver) getNextLines() ([]string, error) {
	incomingBuffer := make([]byte, 1024) // TODO: Determine if this is big enough
	readLen, err := (*sr.serialPort).Read(incomingBuffer)
	if err != nil {
		return []string{}, err
	}

	// TODO: Improve this handling
	bufferString := string(incomingBuffer[:readLen])
	rawSentences := strings.Split(bufferString, "\r\n")

	var notEmptySentences []string
	for _, rawSentence := range rawSentences {
		if rawSentence != "" {
			notEmptySentences = append(notEmptySentences, rawSentence)
		}
	}

	return notEmptySentences, nil
}

func (sr *serialReceiver) parseAndUpdate(sentence string) {
	s, err := nmea.Parse(sentence)
	if err != nil {
		logrus.Warnf("Bad sentence: %v", err)
		return
	}

	switch s.DataType() {
	case nmea.TypeGLL:
		sr.gnssData.UpdateFromGLL(s.(nmea.GLL))
		sr.notifyLocationUpdated()
	case nmea.TypeGGA:
		sr.gnssData.UpdateFromGGA(s.(nmea.GGA))
		sr.notifyLocationUpdated()
	case nmea.TypeVTG:
		sr.gnssData.UpdateFromVTG(s.(nmea.VTG))
		sr.notifyLocationUpdated()
	case nmea.TypeZDA:
		sr.gnssData.UpdateFromZDA(s.(nmea.ZDA))
		sr.notifyLocationUpdated()
	default:
		logrus.Debugf("Unhandled sentence type: %v", s.String())
	}
}

// The GPS continually streams data, but we probably don't need that much data
// This might be called at some interval to keep memory usage low and/or keep the buffer from filling up
func (sr *serialReceiver) resetBuffer() {
	_ = (*sr.serialPort).ResetOutputBuffer()
}
