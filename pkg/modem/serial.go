package modem

import (
	"errors"
	"go.bug.st/serial.v1"
)

type Device struct {
	SerialPort *serial.Port
}

func GetPortNames() (error, []string) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return err, nil
	}
	if len(ports) == 0 {
		return errors.New("no serial ports found"), nil
	}
	return nil, ports
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
