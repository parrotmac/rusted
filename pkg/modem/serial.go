package modem

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"
)

const AT_CMD_AT = "AT\r"
const AT_CMD_GET_MODEL = "AT+CGMM\r"

type Device struct {
	SerialPort *serial.Port
}

func (d *Device) _readSelect(respChannel chan []byte) {
	readBuffer := make([]byte, 256)
	readLen, err := (*d.SerialPort).Read(readBuffer)
	if err != nil {
		logrus.Warnf("Read err: %v", err)
	}

	respChannel <- readBuffer[:readLen]
}

func (d *Device) readDataWithDeadline(deadlineDelta time.Duration) ([]byte, error) {
	var readData []byte
	respChannel := make(chan []byte)
	go d._readSelect(respChannel)
	//defer func() {
	//	close(respChannel)
	//}()
	startTime := time.Now()
	for {
		if time.Now().Sub(startTime) > deadlineDelta {

			return []byte{}, errors.New("reached read deadline")
		}

		select {
		case readData = <-respChannel:
		default:
		}

		if len(readData) > 0 {
			return readData, nil
		}

		// Sleep 1ms
		time.Sleep(time.Millisecond)
	}
}

func GetPortNames() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("no serial ports found")
	}
	return ports, nil
}

func OpenPort(portName string) (error, *serial.Port) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return err, nil
	}
	return nil, &port
}

func (d *Device) writeData(command string) error {
	sp := *d.SerialPort
	writeLen, err := sp.Write([]byte(command))
	if err != nil {
		return err
	}
	logrus.Debugf("Wrote %d/%d bytes", writeLen, len(command))
	return nil
}

func (d *Device) readData() (string, error) {
	sp := *d.SerialPort
	readBuffer := make([]byte, 256)
	readLen, err := sp.Read(readBuffer)
	if err != nil {
		return "", err
	}
	rawStr := string(readBuffer[:readLen])

	strippedString := strings.Trim(rawStr, "\r\n ")

	return strippedString, nil
}

func (d *Device) SendModemCommandWithDeadline(plainCommand string, deadlineDelta time.Duration) (string, error) {
	terminatedCommand := fmt.Sprintf("%s\r", plainCommand)
	writeErr := d.writeData(terminatedCommand)
	if writeErr != nil {
		logrus.Warnf("got error: %v", writeErr)
		return "", writeErr
	}

	for i := 0; i < 5; i++ {
		readData, err := d.readDataWithDeadline(time.Second * 1)
		if err != nil {
			logrus.Errorf("Unable to read data: %v", err)
			continue
		}

		strResp := strings.Trim(string(readData), "\r\n ")

		logrus.Debugf("Read: %s", strResp)

		return strResp, nil
	}
	return "", errors.New("no data returned after max attempts and deadlines")
}

func ClosePort(sp *serial.Port) {
	if sp != nil {
		port := *sp
		_ = port.Close()
	}
}

func checkPort(portName string) (*Device, error) {
	openErr, serialPort := OpenPort(portName)
	if openErr != nil {
		logrus.Debugln("[%s] Unable to open: %v", portName, openErr)
		return nil, openErr
	}

	//if serialPort != nil {
	//	defer ClosePort(serialPort)
	//}

	d := &Device{
		SerialPort: serialPort,
	}

	writeErr := d.writeData(AT_CMD_AT)
	if writeErr != nil {
		logrus.Debugln("[%s] Unable to write: %v", portName, writeErr)
		return nil, writeErr
	}

	for i := 0; i < 5; i++ {
		readData, err := d.readDataWithDeadline(time.Second * 1)
		if err != nil {
			logrus.Errorf("[%s] Unable to read data: %v", portName, err)
		} else {
			strResp := strings.Trim(string(readData), "\r\n ")
			logrus.Infof("[%s] Read: %s", portName, strResp)
			if strResp == "OK" {
				return d, nil
			}
		}
	}

	return nil, errors.New("unable to find corresponding serial port")

	//if readData, err := d.readData(); err != nil {
	//	logrus.Warnf("Did not recognize resp: %s", readData)
	//	return nil, err
	//} else {
	//	logrus.Debugf("Read: %v", readData)
	//}
}

func FindLowSpeedHuaweiModemPort() (*Device, error) {
	portCandidates, err := GetPortNames()
	if err != nil {
		return nil, err
	}
	for _, portName := range portCandidates {
		device, err := checkPort(portName)
		if err != nil {
			logrus.Debugf("Err with %s: %v", portName, err)
		} else {
			return device, nil
		}
	}
	return nil, errors.New("unable to find a suitable serial port")
}
