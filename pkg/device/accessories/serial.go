package accessories

import (
	"errors"
	"time"

	"go.bug.st/serial.v1"
)

type SerialAccessory struct {
	portAddress string
	baudRate    int

	serialPort *serial.Port
}

func OpenSerialPort(deviceAddress string, baudRate int) (*SerialAccessory, error) {
	sa := &SerialAccessory{
		portAddress: deviceAddress,
		baudRate:    baudRate,
	}
	err := openPort(sa)
	if err != nil {
		return nil, err
	}
	return sa, nil
}

func openPort(sa *SerialAccessory) error {
	portMode := serial.Mode{
		BaudRate: sa.baudRate,
	}
	port, err := serial.Open(sa.portAddress, &portMode)
	if err != nil {
		return err
	}
	sa.serialPort = &port
	return nil
}

func (sa *SerialAccessory) write(data []byte) (int, error) {
	return (*sa.serialPort).Write(data)
}

func (sa *SerialAccessory) read(dest *[]byte) (int, error) {
	return (*sa.serialPort).Read((*dest))
}

func (sa *SerialAccessory) SendCommand(command []byte) ([]byte, error) {
	writeLen, err := sa.write(command)
	if err != nil {
		// TODO: Optionally retry
		return nil, err
	}

	// TODO: See if this is thrown by internal serial port
	if writeLen != len(command) {
		return nil, errors.New("not all data was written in first attempt")
	}

	buf := make([]byte, 1024)

	readLen, err := sa.read(&buf)

	if err != nil {
		return nil, err
	}

	// TODO: Does this panic if the readLen is oob?
	return buf[:readLen], nil
}

func (sa *SerialAccessory) ReadLine() ([]byte, error) {
	return nil, errors.New("Not yet implemented")
}

func (sa *SerialAccessory) ReadLineTimeout(timeout time.Duration) ([]byte, error) {
	return nil, errors.New("Not yet implemented")
}
