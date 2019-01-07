package gnss

import (
	"errors"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"

	"github.com/parrotmac/rusted/pkg/central"
	"github.com/parrotmac/rusted/pkg/central/entities"
)

type BasicLocationUpdateDelegate func(l entities.BasicLocation)
type AdvancedLocationUpdateDelegate func(l entities.AdvancedLocation)

type serialReceiver struct {
	deviceAddress string
	baudRate      int

	serialPort *serial.Port

	skipValidityChecks bool

	basicUpdateDelegate    BasicLocationUpdateDelegate
	advancedUpdateDelegate AdvancedLocationUpdateDelegate

	// TODO: Break these out
	readValueTimeout time.Duration

	basicLocationLastReadTime time.Time
	basicLocationLast         entities.BasicLocation

	advancedLocationLastReadTime time.Time
	advancedLocationLast         entities.AdvancedLocation
}

func (sr *serialReceiver) Start(ctx *central.Context) error {
	// TODO: Rethink this
	go sr.processMessagesContinually()
	return nil
}

func (sr *serialReceiver) GetBasicLocation(ctx *central.Context) (entities.BasicLocation, error) {
	oldReadTime := sr.basicLocationLastReadTime
	startTime := time.Now()
	for {
		if sr.basicLocationLastReadTime.After(oldReadTime) {
			return sr.basicLocationLast, nil
		}
		if time.Now().After(startTime.Add(sr.readValueTimeout)) {
			return entities.BasicLocation{}, errors.New("timeout waiting for new basic location data")
		}
	}
}

func (sr *serialReceiver) GetDetailedLocation(ctx *central.Context) (entities.AdvancedLocation, error) {
	oldReadTime := sr.advancedLocationLastReadTime
	startTime := time.Now()
	for {
		if sr.advancedLocationLastReadTime.After(oldReadTime) {
			return sr.advancedLocationLast, nil
		}
		if time.Now().After(startTime.Add(sr.readValueTimeout)) {
			return entities.AdvancedLocation{}, errors.New("timeout waiting for new advanced location data")
		}
	}
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
		deviceAddress:      deviceAddress,
		baudRate:           baudRate,
		skipValidityChecks: false,
		readValueTimeout:   time.Duration(time.Second * 2),
	}

	port, err := sr.openPort()
	if err != nil {
		return nil, err
	}
	sr.serialPort = port

	return sr, nil
}

func (sr *serialReceiver) SetValitiyChecking(enforceValitity bool) {
	sr.skipValidityChecks = !enforceValitity
}

func (sr *serialReceiver) SetBasicUpdateDelegate(dataUpdatedDelegate BasicLocationUpdateDelegate) {
	sr.basicUpdateDelegate = dataUpdatedDelegate
}

func (sr *serialReceiver) SetAdvancedUpdateDelegate(dataUpdatedDelegate AdvancedLocationUpdateDelegate) {
	sr.advancedUpdateDelegate = dataUpdatedDelegate
}

func (sr *serialReceiver) notifyGLLUpdate(gll nmea.GLL) {

	basicLoc := entities.NewBasicLocationFromGLL(gll)
	logrus.Debugf("Notifying of update to gll: %v", gll)

	// TODO: Separate these
	sr.basicLocationLast = basicLoc
	sr.basicLocationLastReadTime = time.Now()

	// TODO: Reevaluate these
	if sr.basicUpdateDelegate != nil {
		if gll.Validity == "A" || sr.skipValidityChecks {
			sr.basicUpdateDelegate(basicLoc)
			return
		}
		logrus.Debugf("Not notifying of update to gll: %v because of bad validity", gll)
	} else {
		logrus.Debugln("basicUpdateDelegate isn't set")
	}
}

func (sr *serialReceiver) notifyGGAUpdate(gga nmea.GGA) {

	logrus.Debugf("Notifying of update to gga: %v", gga)
	advLoc := entities.NewAdvancedLocationFromGGA(gga)

	// TODO: Separate these
	sr.advancedLocationLast = advLoc
	sr.advancedLocationLastReadTime = time.Now()

	// TODO: Reevaluate this
	if sr.advancedUpdateDelegate != nil {
		if gga.FixQuality != "0" || sr.skipValidityChecks {
			sr.advancedUpdateDelegate(advLoc)
			return
		}
		logrus.Debugf("Not notifying of update to gga: %v because of bad fix quality", gga)
	} else {
		logrus.Debugln("advancedUpdateDelegate isn't set")
	}
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

	sentencePrefix := s.Prefix()

	logrus.Debugln("sentencePrefix", sentencePrefix)

	switch sentencePrefix {
	// TODO; Make agnostic to talker type
	case "GNGLL":
		sr.notifyGLLUpdate(s.(nmea.GLL))
		break
	case "GNGGA":
		sr.notifyGGAUpdate(s.(nmea.GGA))
		break
	default:
		logrus.Debugf("Unhandled sentence type: %v", s.String())
	}
}

// The GPS continually streams data, but we probably don't need that much data
// This might be called at some interval to keep memory usage low and/or keep the buffer from filling up
func (sr *serialReceiver) resetBuffer() {
	_ = (*sr.serialPort).ResetOutputBuffer()
}
