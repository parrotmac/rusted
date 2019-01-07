package accessories

import "time"

type StreamAccessory interface {
	SendCommand(command []byte) ([]byte, error)
	ReadLine() ([]byte, error)
	ReadLineTimeout(timeout time.Duration) ([]byte, error)
}
