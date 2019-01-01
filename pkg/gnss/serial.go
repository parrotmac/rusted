package gnss

import (
	"github.com/adrianmo/go-nmea"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"
	"strings"
)

type BasicLocationUpdateDelegate func(l BasicLocation)
type AdvancedLocationUpdateDelegate func(l AdvancedLocation)

type serialReceiver struct {
	deviceAddress string
	baudRate      int

	serialPort *serial.Port

	skipValidityChecks bool

	basicUpdateDelegate    BasicLocationUpdateDelegate
	advancedUpdateDelegate AdvancedLocationUpdateDelegate
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
	if sr.basicUpdateDelegate != nil {
		if gll.Validity == "A" || sr.skipValidityChecks {
			logrus.Debugf("Notifying of update to gll: %v", gll)
			sr.basicUpdateDelegate(newBasicLocationFromGLL(gll))
			return
		}
		logrus.Debugf("Not notifying of update to gll: %v because of bad validity", gll)
	} else {
		logrus.Debugln("basicUpdateDelegate isn't set")
	}
}

func (sr *serialReceiver) notifyGGAUpdate(gga nmea.GGA) {
	if sr.advancedUpdateDelegate != nil {
		if gga.FixQuality != "0" || sr.skipValidityChecks {
			logrus.Debugf("Notifying of update to gga: %v", gga)
			sr.advancedUpdateDelegate(newAdvancedLocationFromGGA(gga))
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
